package agendaCulturalPorto

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	m "github.com/oleksiy-os/porto-events/internal/model"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SourceAgendaculturalPorto struct {
	m.Source
}

func (s *SourceAgendaculturalPorto) LoadEvents(u *url.URL) []m.Event {
	var (
		ev     m.Event
		events []m.Event
	)

	res, err := http.Get(u.String())
	if err != nil {
		log.Error(err, u.String())
		return events
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	// Load the HTML document
	doc, err := m.LoadContent(res)
	if err != nil {
		log.Errorln("error load content ", err)
		return events
	}

	// Find the review items
	doc.Find("article.mec-event-article").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".mec-event-title a")
		if title.Length() == 0 {
			log.Error("parse event| not found title", u.String())
			return
		}
		ev = m.Event{
			Title: title.Text(),
		}
		ev.ID, _ = title.Attr("data-event-id")
		ev.Url, _ = title.Attr("href")
		srcSet, ok := s.Find(".mec-event-image img").Attr("data-lazy-srcset")
		if !ok {
			log.Error("image not found", ev.Url)
		} else {
			ev.Image, err = image(srcSet)
			if err != nil {
				log.Error(err, ev.Url)
			}
		}

		log.Debugln("visiting ev page for more data collect")
		eventPageUrl, err := url.ParseRequestURI(ev.Url)
		if err != nil {
			log.Error("failed event url", ev.Url)
		}
		eventPageUrl.Scheme = u.Scheme // need for proper tests work
		eventPageUrl.Host = u.Host     // need for proper tests work

		res, err := http.Get(eventPageUrl.String())
		if err != nil {
			log.Error(err)
			return
		}

		//goland:noinspection GoUnhandledErrorResult
		defer res.Body.Close()
		// Load the HTML document
		d, err := m.LoadContent(res)
		if err != nil {
			log.Errorln("error load content ", err)
			return
		}
		el := d.Find(".mec-single-event").First()
		if el.Length() == 0 {
			log.Error("parse event page| not found wrap", eventPageUrl)
			return
		}
		ev.Place = el.Find(".mec-single-event-location .author").Text()
		ev.Location = el.Find(".mec-single-event-location .mec-address").Text()
		ev.Description = m.StripAllHtml.Sanitize(el.Find(".mec-single-event-description p").Text())
		ev.Time = el.Find(".mec-single-event-time .mec-events-abbr").Text()
		ev.DateText = monthPtToEn(el.Find(".mec-single-event-date .mec-events-abbr .mec-start-date-label").Text())
		ev.Timestamp, err = timestamp(ev.DateText, ev.Time)
		if err != nil {
			log.Error("date parse", err, ev.DateText, ev.Time, ev.Url)
			return
		}

		log.Print("ev Description ", ev.Description)
		events = append(events, ev)
	})

	return events
}

func image(srcSet string) (string, error) {
	lastIndex := strings.Index(srcSet, " 300w")
	if lastIndex == -1 {
		return "", errors.New("not found 300px img src ")
	}
	firstIndex := strings.LastIndex(srcSet[:lastIndex], "://")
	if lastIndex == -1 {
		return "", errors.New("not found 300px img src ")
	}
	firstIndex = firstIndex - 5 // len("https")

	return strings.TrimLeft(srcSet[firstIndex:lastIndex], " "), nil
}

func New(sourceConfig m.Source) *SourceAgendaculturalPorto {
	return &SourceAgendaculturalPorto{
		m.Source{
			Name: sourceConfig.Name,
			Url:  sourceConfig.Url,
		},
	}
}

func monthPtToEn(datePt string) string {
	var month = map[string]string{
		"Fev": "Feb",
		"Abr": "Apr",
		"Mai": "May",
		"Ago": "Aug",
		"Set": "Sep",
		"Out": "Oct",
		"Dez": "Dec",
	}

	for pt, en := range month {
		if strings.Contains(datePt, pt) {
			return strings.Replace(datePt, pt, en, 1)
		}
	}
	return datePt
}

// timestamp generate
//
// date: "06 Jan 2023"
//
// timeTxt: "21:00 - 23:30"
func timestamp(date string, timeTxt string) (time.Time, error) {
	layoutInput := "02 Jan 2006 15:04"

	i := strings.Index(timeTxt, " - ")
	if i != -1 {
		timeTxt = timeTxt[:i]
	}

	start, err := time.Parse(layoutInput, date+" "+timeTxt)
	if err != nil {
		return time.Now(), err
	}

	if start.Before(time.Now()) {
		return time.Now().Truncate(time.Hour).In(time.FixedZone("WET", 0)), nil
	}

	return start, nil
}
