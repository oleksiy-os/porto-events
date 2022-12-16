package porto

import (
	"encoding/json"
	"fmt"
	m "github.com/oleksiy-os/porto-events/internal/model"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (
	SourcePorto struct {
		Name string
		Url  string
	}

	Dates struct {
		Start     string `json:"start"`
		End       string `json:"end"`
		Repeating []struct {
			Label string `json:"label"`
		}
	}

	EventSource struct {
		Id        string   `json:"id"`
		Url       string   `json:"url"`
		FullUrl   string   `json:"fullUrl"`
		Title     string   `json:"title"`
		Dates     [1]Dates `json:"dates"`
		Thumbnail struct {
			Small struct {
				Url string `json:"url"`
			} `json:"small"`
		} `json:"thumbnail"`
		Locations []struct {
			Location struct {
				Locality  string
				Address   string
				Latitude  float64
				Longitude float64
			}
		}
	}

	EventList struct {
		PageByUrl struct {
			Events struct {
				Items []EventSource `json:"items"`
			} `json:"events"`
		} `json:"pageByUrl"`
	}
)

func (s *SourcePorto) LoadEvents(u *url.URL) []m.Event {
	var (
		event  m.Event
		events []m.Event
	)

	eventsData, err := getFromApi(u)
	if err != nil {
		log.Errorln("get from api| ", u.String(), err)
		return events
	}

	for _, ev := range *eventsData {
		event = m.Event{
			ID:          m.StripAllHtml.Sanitize(ev.Id),
			Title:       m.StripAllHtml.Sanitize(ev.Title),
			Description: description(u.Scheme+"://"+u.Host+u.Path, ev.Url),
			Url:         ev.FullUrl,
			Image:       ev.Thumbnail.Small.Url,
			Days:        parseDays(ev.Dates[0]),
			Place:       parsePlace(ev),
			LocationMap: fmt.Sprintf(
				"https://www.google.com/maps/search/?api=1&query=%f,%f",
				ev.Locations[0].Location.Latitude,
				ev.Locations[0].Location.Longitude,
			),
			Timestamp: getTime(ev.Dates[0]),
		}

		event.DateText, event.Time = parseDate(ev.Dates[0])

		events = append(events, event)
	}

	return events
}

func getTime(dates Dates) time.Time {
	layoutInput := "2006-01-02 15:04:05" // input format in json 2022-05-12 10:00:00

	start, err := time.Parse(layoutInput, dates.Start)
	if err != nil {
		log.Error("date parse", err, dates.Start)
	}

	if start.Before(time.Now()) {
		return time.Now().Truncate(time.Hour).In(time.FixedZone("WET", 0))
	}

	return start
}

func parsePlace(event EventSource) string {
	loc := event.Locations[0].Location
	if loc.Address != "" {
		return loc.Locality +
			" - " +
			loc.Address
	}

	return m.StripAllHtml.Sanitize(loc.Locality)
}

func getFromApi(apiUrl *url.URL) (eventsSource *[]EventSource, err error) {
	var data *EventList

	values := apiUrl.Query()
	values.Set("startDate", time.Now().Format("2006-01-02"))
	apiUrl.RawQuery = values.Encode()

	res, err := http.Get(apiUrl.String())
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Error(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("read body", err, apiUrl.String())
		return nil, err
	}

	if err = json.Unmarshal(body, &data); err != nil {
		log.Error("unmarshal - ", err, apiUrl.String())
		return nil, err
	}

	return &data.PageByUrl.Events.Items, nil
}

func description(apiUrl string, eventPagePath string) string {
	var descr struct {
		PageByUrl struct {
			Body []struct {
				Value string `json:"value"`
			} `json:"body"`
		}
	}

	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		log.Error("parse apiUrl", err, apiUrl)
	}

	values := u.Query()
	values.Add("queryName", "PageByUrl")
	values.Add("urlPath", eventPagePath)
	u.RawQuery = values.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return ""
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Error(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("read body", err, u.String())
		return ""
	}

	err = json.Unmarshal(body, &descr)
	if err != nil || descr.PageByUrl.Body == nil {
		log.Error("unmarshal description ", err, u.String())
		return ""
	}

	return m.StripAllHtml.Sanitize(descr.PageByUrl.Body[0].Value)
}

func parseDays(dates Dates) string {
	var days = ""
	for _, item := range dates.Repeating {
		days += item.Label + ", "
	}

	return m.StripAllHtml.Sanitize(strings.TrimRight(days, ", "))
}

func parseDate(dates Dates) (DateText string, Time string) {
	layoutInput := "2006-01-02 15:04:05" // input format in json 2022-05-12 10:00:00
	layoutOutDate := "Jan 02th, 2006"
	layoutOutTime := "15:04"

	start, err := time.Parse(layoutInput, dates.Start)
	if err != nil {
		log.Error("date parse", err, dates.Start)
	}
	end, err := time.Parse(layoutInput, dates.End)
	if err != nil {
		log.Error("date parse", err, dates.Start)
	}

	return start.Format(layoutOutDate) + " - " + end.Format(layoutOutDate),
		start.Format(layoutOutTime) + " - " + end.Format(layoutOutTime)
}

func New(sourceConfig m.Source) *SourcePorto {
	return &SourcePorto{
		Name: sourceConfig.Name,
		Url:  sourceConfig.Url,
	}
}
