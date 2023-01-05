package agendaCulturalPorto

import (
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

func TestSourceAgendaculturalPorto_LoadEvents(t *testing.T) {
	svr := httptest.NewServer(requestHandler(t))

	defer svr.Close()
	source := New(model.Source{
		Name: "test",
		Url:  svr.URL,
	})

	u, _ := url.Parse(source.Url)
	evs := source.LoadEvents(u)

	tests := []struct {
		name string
		want model.Event
		got  []model.Event
	}{
		{
			name: "ok",
			got:  evs,
			want: model.Event{
				ID:    "35215",
				Title: "Orfélia em estreia ao vivo no Maus Hábitos",
				Url:   "/orfelia-em-estreia-ao-vivo-no-maus-habitos",
				Description: `O PROJETO MUSICAL LUSO-BRASILEIRO APRESENTA O ÁLBUM DE ESTREIA “TUDO O QUE MOVE” – 6 DE JANEIRO – MAUS HÁBITOSO projeto luso-brasileiro formado por Antera e Filipe Mattos, vão 
atuar pela primeira vez na cidade do Porto: os Orfélia vão apresentar o 
disco de estreia “Tudo o que Move” no dia 6 de Janeiro, no Maus Hábitos,
 às 21h30.O disco de estreia mistura ritmos tradicionais com o mistério 
envolvente do psicadelismo. Filha do tropicalismo, expõe e celebra a 
música multicultural entre dois povos irmãos. O título “Tudo que Move” é
 inspirado na alma da música de Gilberto Gil “Aqui e Agora”, que marca 
um momento da vida dos artistas de muita transformação, revelador da 
força do espírito.Antera, natural de Lagos (Algarve), estudou piano clássico e canto e 
formou-se em Artes Performativas; Filipe estudou guitarra em 
Conservatórios de música e é natural de Florianópolis, uma cidade no sul
 do Brasil com uma forte influência portuguesa. Conheceram-se num palco 
em Berlim e a partir desse encontro nasceu “Orfélia”.As suas influências musicais vão desde Chico Buarque, Jacques Brel, 
The Beatles, até Amália Rodrigues, Caetano Veloso, entre outros mestres 
intemporais.Em 2019 o duo lançou o EP “Retratos Temporais” e singles que 
receberam destaque diversos em Meios de comunicação. Em especial o 
single “Lagos”, um dos temas vencedores do concurso Inéditos Vodafone, 
promovido pela Vodafone e Sony Music Portugal em 2020.No dia 6 de Janeiro os Orfélia vão apresentar-se com banda completa 
que conta com Antera na voz e sintetizadores, Filipe Mattos na guitarra,
 André Morais no baixo, Sebastião Bergmann na bateria, Lana Gasparotti 
nas teclas e Zé Cruz na percussão.Orfélia em estreia ao vivo nos Maus Hábitos`,
				Image:     "https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-150x150.jpg",
				Place:     "Maus Hábitos - Espaço de Intervenção Cultural",
				Location:  "R. de Passos Manuel 178 4º Piso, 4000-382 Porto",
				DateText:  "06 Jan 2024",
				Time:      "21:00 - 23:30",
				Timestamp: time.Date(2024, time.January, 6, 21, 0, 0, 0, time.FixedZone("WET", 0)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Len(t, tt.got, 1)
			assert.Equal(t, tt.want.ID, tt.got[0].ID)
			assert.Equal(t, tt.want.Title, tt.got[0].Title)
			assert.Equal(t, tt.want.Url, tt.got[0].Url)
			assert.Equal(t, tt.want.Image, tt.got[0].Image)
			assert.Equal(t, tt.want.Description, tt.got[0].Description)
			assert.Equal(t, tt.want.Place, tt.got[0].Place)
			assert.Equal(t, tt.want.Location, tt.got[0].Location)
			assert.Equal(t, tt.want.DateText, tt.got[0].DateText)
			assert.Equal(t, tt.want.Time, tt.got[0].Time)
			assert.Truef(t, tt.want.Timestamp.Equal(tt.got[0].Timestamp), "want: %s, got: %s", tt.want.Timestamp, tt.got[0].Timestamp)
		})
	}
}

func requestHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := "tests/"
		filePath := path + "eventsList.html"
		eventPage := "orfelia-em-estreia-ao-vivo-no-maus-habitos"

		if strings.Contains(r.URL.String(), eventPage) {
			filePath = path + eventPage + ".html"
		}

		file, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatal("file not found|", err, r.RequestURI)
		}

		//w.Header().Set("Content-Type", "application/json")
		if _, err = w.Write(file); err != nil {
			t.Fatal("write file|", err)
		}
	}
}

func Test_timestamp(t *testing.T) {
	type args struct {
		date    string
		timeTxt string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "ok, date in future",
			args: args{
				date:    "06 Jan 2024",
				timeTxt: "21:00 - 23:30",
			},
			want:    time.Date(2024, time.January, 6, 21, 0, 0, 0, time.FixedZone("WET", 0)),
			wantErr: false,
		},
		{
			name: "ok, date in past, use timeNow",
			args: args{
				date:    "02 Jan 2022",
				timeTxt: "21:00 - 23:30",
			},
			want:    time.Now().Truncate(time.Hour).In(time.FixedZone("WET", 0)),
			wantErr: false,
		},
		{
			name: "date parse err",
			args: args{
				date:    "Jan 2022",
				timeTxt: "21:00 - 23:30",
			},
			wantErr: true,
		},
		{
			name: "time only start",
			args: args{
				date:    "06 Jan 2023",
				timeTxt: "21:00",
			},
			wantErr: false,
			want:    time.Date(2023, time.January, 6, 21, 0, 0, 0, time.FixedZone("WET", 0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := timestamp(tt.args.date, tt.args.timeTxt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Truef(t, tt.want.Equal(got), "wrong date, want: %s, got: %s", tt.want, got)
			}
		})
	}
}

func Test_image(t *testing.T) {
	type args struct {
		srcSet string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				srcSet: `https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-300x300.jpg 300w, https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-1024x1024.jpg 1024w, https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-150x150.jpg 150w, https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-768x768.jpg 768w, https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia.jpg 1200w" data-lazy-sizes="(max-width: 300px) 100vw, 300px" data-lazy-src="https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-300x300.jpg`,
			},
			want:    "https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-150x150.jpg",
			wantErr: false,
		},
		{
			name: "no 150w",
			args: args{
				srcSet: `https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-300x300.jpg 300w, https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-1024x1024.jpg 1024w, https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-768x768.jpg 768w, https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia.jpg 1200w" data-lazy-sizes="(max-width: 300px) 100vw, 300px" data-lazy-src="https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-300x300.jpg`,
			},
			wantErr: true,
		},
		{
			name: "http img url",
			args: args{
				srcSet: `https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-300x300.jpg 300w, https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-1024x1024.jpg 1024w, http://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-150x150.jpg 150w, https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-768x768.jpg 768w, https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia.jpg 1200w" data-lazy-sizes="(max-width: 300px) 100vw, 300px" data-lazy-src="https://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-300x300.jpg`,
			},
			want:    "http://agendaculturalporto.org/wp-content/uploads/2022/12/Orfelia-150x150.jpg",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := image(tt.args.srcSet)
			if tt.wantErr {
				assert.Error(t, err, err)
				return
			}

			assert.Equal(t, tt.want, got, "wrong img url")
		})
	}
}
