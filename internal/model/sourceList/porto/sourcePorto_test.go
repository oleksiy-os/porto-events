package porto

import (
	"fmt"
	"github.com/oleksiy-os/porto-events/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_getTime(t *testing.T) {
	tests := []struct {
		name  string
		input Dates
		want  time.Time
	}{
		{
			name:  "date before time.Now",
			input: Dates{Start: "2022-05-12 10:00:00"},
			want:  time.Now().Truncate(time.Hour).In(time.FixedZone("WET", 0)),
		},
		{
			name:  "ok",
			input: Dates{Start: "2024-05-12 10:00:00"},
			want:  time.Date(2024, time.May, 12, 10, 0, 0, 0, time.FixedZone("WET", 0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			get := getTime(tt.input)
			assert.True(t, tt.want.Equal(get), tt.name)
		})
	}
}

func Test_getFromApi(t *testing.T) {
	tests := []struct {
		name   string
		apiUrl func() *url.URL
		wantOk bool
	}{
		{
			name: "ok",
			apiUrl: func() *url.URL {
				urlQuery := fmt.Sprintf(
					"https://www.porto.pt/api/graphql?queryName=PageByUrl&urlPath=/en/events/&startDate=%s",
					time.Now().Format("2006-01-02"),
				)
				u, _ := url.ParseRequestURI(urlQuery)

				return u
			},
			wantOk: true,
		},
		{
			name: "wrong url",
			apiUrl: func() *url.URL {
				u, _ := url.Parse("wrongUrl")

				return u
			},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEventsSource, err := getFromApi(tt.apiUrl())

			if tt.wantOk == false {
				assert.Error(t, err, tt.apiUrl().String())
				return
			}

			if !assert.NoError(t, err, err) {
				return
			}

			if assert.True(t, len(*gotEventsSource) > 0, "events count %v", len(*gotEventsSource)) {
				for _, e := range *gotEventsSource {
					assert.NotEmpty(t, e.Id, "id")
					assert.NotEmpty(t, e.Url, "url")
					assert.NotEmpty(t, e.FullUrl, "full url")
					assert.NotEmpty(t, e.Title, "title")
					assert.NotEmpty(t, e.Dates, "dates")
					assert.NotEmpty(t, e.Thumbnail.Small.Url, "image")
					assert.NotEmpty(t, e.Locations[0].Location.Locality, "Locality")
					assert.NotEmpty(t, e.Locations[0].Location.Latitude, "Latitude")
					assert.NotEmpty(t, e.Locations[0].Location.Longitude, "Longitude")
					if t.Failed() {
						t.Fatalf("apiUrl: %s, event: %+v", tt.apiUrl().String(), e)
					}
				}
			}
		})
	}
}

func TestSourcePorto_LoadEvents(t *testing.T) {
	timeNow := time.Now().Truncate(time.Hour).In(time.FixedZone("WET", 0))
	timeLayout := "2006-01-02 15:04:05"
	tests := []struct {
		name string
		want model.Event
	}{
		{
			name: "ok 36011",
			want: model.Event{
				ID:          "36011",
				Url:         "https://www.porto.pt/en/event/exhibition-este-mundo-nao-nos-pertence/",
				Title:       "Exhibition | “Este mundo não nos pertence”",
				Description: `The temporary exhibition “Este mundo não nos pertence” (&#34;This world does not belong to us&#34;) marks the fourth great moment of artistic and cultural dissemination of Espaço João Espregueira Mendes (EJEM), at Museu Futebol Clube do Porto, with national and foreign authors.From the private collection of Isabel Mota and Fernando Pereira, painting, sculpture, drawing and photography come together in an interesting approach to contemporary art from the end of the 20th century to the present.Open to the public on May 12, the exhibition has been extended and can continue to be visited until December 31. Entry is free.Complete information at Futebol Clube do Porto website.` + "\n",
				Place:       "Museu FC Porto - Espaço João Espregueira Mendes",
				Image:       "https://www.porto.pt/_next/image?url=https%3A%2F%2Fmedia.porto.pt%2Foriginal_images%2FIsabel_Mota_Fernando_Pereira_exposicao_Este_mundo_nao_nos_pertence_01.JPG&w=350&q=85",
				LocationMap: "https://www.google.com/maps/search/?api=1&query=41.161367,-8.583016",
				Timestamp:   timeNow,
				DateText:    "May 12th, 2022 - Dec 31th, 2022",
				Days:        "mon, tue, wed, thu, fri, sat, sun",
				Time:        "10:00 - 18:00",
			},
		},
		{
			name: "ok 36012",
			want: model.Event{
				ID:          "36012",
				Url:         "https://www.porto.pt/en/event/exhibition-fictional-grounds/",
				Title:       "Exhibition | “Fictional Grounds”",
				Description: `Soil simulations of an imagined territory through which one can look for traces of minerals with energetic potential and present soil samples from different origins with varied compositions that are mounted in two-dimensional planes – it is this fictional reality that can be watched closely in the new exhibition “Fictional Grounds” by the artistic collective berru.Created in Porto, the collective won the Sonae Media Art 2019 award, and has already exhibited and was responsible for installations at institutions such as Calouste Gulbenkian Foundation, BoCA Biennial of Contemporary Arts, and The Old Truman Brewey (London).Viewed critically and driven by the urgency of the current ecological catastrophe, the exhibition establishes a subtle relationship with the world of earthworks by pioneering Land Art artists such as Robert Smithson, Richard Long or the famous exhibition by Walter de Maria when in 1977 he filled a gallery in New York with 140 tons of earth.Entry is free.` + "\n",
				Place:       "Escola das Artes da Universidade Católica Portuguesa",
				Image:       "https://www.porto.pt/_next/image?url=https%3A%2F%2Fmedia.porto.pt%2Foriginal_images%2FDR_Fictional_Grounds_exposicao_coletivo_berru.jpg&w=350&q=85",
				LocationMap: "https://www.google.com/maps/search/?api=1&query=41.154207,-8.672795",
				Timestamp:   timeNow,
				DateText:    "Oct 20th, 2022 - Feb 17th, 2023",
				Days:        "mon, tue, wed, thu, fri, sat, sun",
				Time:        "10:00 - 19:00",
			},
		},
		{
			name: "ok 35974",
			want: model.Event{
				ID:          "35974",
				Url:         "https://www.porto.pt/en/event/exhibition-walking-art-maps/",
				Title:       "Exhibition | Walking Art Maps",
				Description: `The Walking Art Maps – #asbelasarteseacidade exhibition, especially dedicated to international students, within the scope of the 35th anniversary of ERASMUS+, opens this Wednesday at the Exhibition Pavilion of the Faculty of Fine Arts of the U. Porto.The exhibition brings together works from the FBAUP collection that are associated with works installed in public and private spaces, easily accessible. Six routes are proposed, where works by some of the artists, architects and designers of the Belas Artes do Porto are identified, which punctuate and characterize the city in various ways.Walking Art Maps unfolds between the FBAUP Exhibition Pavilion and an online platform. Complete information at FBAUP website. ` + "\n",
				Place:       "Porto - Faculdade de Belas Artes",
				Image:       "https://www.porto.pt/_next/image?url=https%3A%2F%2Fmedia.porto.pt%2Foriginal_images%2FDR_pavilhao_de_exposicoes_FBAUP.jpg&w=350&q=85",
				LocationMap: "https://www.google.com/maps/search/?api=1&query=41.145647,-8.600677",
				Timestamp:   timeNow,
				DateText:    "Oct 26th, 2022 - Jan 14th, 2023",
				Days:        "mon, tue, wed, thu, fri, sat, sun",
				Time:        "17:30 - 18:00",
			},
		},
		{
			name: "ok 35973",
			want: model.Event{
				ID:          "35973",
				Url:         "https://www.porto.pt/en/event/exhibition-so-what/",
				Title:       "Exhibition | So What",
				Description: `Three architects - Diogo Aguiar, Dulcineia Santos and Nuno Melo Sousa - proposed to illuminate the conceptual bases of their three architectural discourses in the exhibition entitled &#34;SO WHAT&#34;,  evocative of Duke Ellington&#39;s restlessness and continual experimentation.The opening session will include a three-way conversation, moderated by Hélder Casal Ribeiro, curator of the exhibition.&#34;SO WHAT&#34; will feature a parallel set of conferences, guided tours and Educational Service activities, upon registration.  Complete information at Casa das Artes website.` + "\n",
				Place:       "Porto - Casa das Artes",
				Image:       "https://www.porto.pt/_next/image?url=https%3A%2F%2Fmedia.porto.pt%2Foriginal_images%2Fc48fd81f0c34-Por_do_Sol_nas_artes_Casa_das_Artes.jpg&w=350&q=85",
				LocationMap: "https://www.google.com/maps/search/?api=1&query=41.156465,-8.643391",
				Timestamp:   timeNow,
				DateText:    "Nov 05th, 2022 - Dec 23th, 2022",
				Days:        "mon, tue, wed, thu, fri, sat, sun",
				Time:        "15:30 - 19:00",
			},
		},
		{
			name: "ok 36013",
			want: model.Event{
				ID:          "36013",
				Url:         "https://www.porto.pt/en/event/show-impossible-by-luis-de-matos/",
				Title:       "Show | IMPOSSIBLE, by Luís de Matos",
				Description: `“Luís de Matos IMPOSSIBLE Live” ends its national tour at Coliseu do Porto with shows on January 13th and 14th.The most awarded Portuguese magician, distinguished three times by the Academy of Magical Arts in Hollywood, and the youngest in history to receive the Devant Award, from The Magic Circle, brings a new journey through the world of illusion where the impossible becomes reality and the limits of imagination are challenged at every moment.Luís de Matos will have at his side four of the greatest magicians in the world today: from the United States of America, Dan Sperry, from Spain, Javier Botía, from France, Norbert Ferré, and from South Korea, Yu Hojin.Complete information at Coliseu do Porto website.` + "\n",
				Place:       "Porto - Coliseu Porto Ageas",
				Image:       "https://www.porto.pt/_next/image?url=https%3A%2F%2Fmedia.porto.pt%2Foriginal_images%2FDR_Luis_de_Matos_coliseu.jpg&w=350&q=85",
				LocationMap: "https://www.google.com/maps/search/?api=1&query=41.146992,-8.605417",
				Timestamp:   time.Date(2023, time.January, 13, 21, 0, 0, 0, time.FixedZone("WET", 0)),
				DateText:    "Jan 13th, 2023 - Jan 14th, 2023",
				Days:        "fri, sat",
				Time:        "21:00 - 18:00",
			},
		},
	}

	svr := httptest.NewServer(requestHandler(t))

	defer svr.Close()
	source := New(model.Source{
		Name: "test",
		Url:  svr.URL,
	})

	u, _ := url.Parse(source.Url)

	evs := source.LoadEvents(u)

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want.ID, evs[i].ID, "ID")
			assert.Equal(t, tt.want.Url, evs[i].Url, "Url")
			assert.Equal(t, tt.want.Title, evs[i].Title, "Title")
			assert.Equal(t, tt.want.Description, evs[i].Description, "Description")
			assert.Equal(t, tt.want.Place, evs[i].Place, "Place")
			assert.Equal(t, tt.want.Image, evs[i].Image, "Image")
			assert.Equal(t, tt.want.LocationMap, evs[i].LocationMap, "LocationMap")
			assert.Equal(t, tt.want.Timestamp.Format(timeLayout), evs[i].Timestamp.Format(timeLayout))
			assert.Equal(t, tt.want.DateText, evs[i].DateText, "DateText")
			assert.Equal(t, tt.want.Days, evs[i].Days, "Days")
			assert.Equal(t, tt.want.Time, evs[i].Time, "Time")
		})
	}
}

func requestHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventViewPage := "/en/event/"
		filePath := "tests/test_events_list.json"

		urlPath := r.URL.Query().Get("urlPath")
		if urlPath != "" && strings.Contains(urlPath, eventViewPage) { // event view page
			urlPath = strings.TrimRight(r.URL.Query().Get("urlPath"), "/")
			filePath = "tests/" + strings.Replace(urlPath, eventViewPage, "", 1) + ".json"
		}

		file, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatal("file not found|", err, r.RequestURI)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err = w.Write(file); err != nil {
			t.Fatal("write file|", err)
		}
	}
}
