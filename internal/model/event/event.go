package event

import (
	"github.com/oleksiy-os/porto-events/internal/model"
	"github.com/oleksiy-os/porto-events/internal/model/sourceList/agendaCulturalPorto"
	"github.com/oleksiy-os/porto-events/internal/model/sourceList/porto"
	"github.com/oleksiy-os/porto-events/internal/model/sourceList/testing"
	log "github.com/sirupsen/logrus"
	"net/url"
	"sort"
)

type (
	sourceInterface interface {
		LoadEvents(url *url.URL) []model.Event
	}
)

// Collect (web scrap OR get from API) events from sources
func Collect(sources []model.Source) *[]model.Event {
	var (
		events           []model.Event
		eventsCollection []model.Event
		src              sourceInterface
	)

	for _, item := range sources {
		log.Debugln("getting events from source:", item.Name, item.Url)

		u, err := url.ParseRequestURI(item.Url)
		if err != nil {
			log.Error("wrong source url", item)
			continue
		}

		switch item.Name {
		case "porto":
			src = porto.New(item)
		case "agendaculturalporto":
			src = agendaCulturalPorto.New(item)
		case "testing":
			src = testing.New() // only for tests purpose

		// ** will be added soon **
		// *********
		//case "teatromunicipaldoporto":
		//	src = sourceList.Teatromunicipaldoporto(item)

		default:
			log.Errorln("undefined source name from source list file", item)
			continue
		}

		if events = src.LoadEvents(u); len(events) == 0 {
			log.Println("No events, strange...", item.Name)
			continue
		}

		log.WithFields(log.Fields{"source": item.Name}).Debugln("Collected events")
		eventsCollection = append(eventsCollection, events...)
	}

	sort.Slice(eventsCollection, func(i, j int) bool {
		return eventsCollection[i].Timestamp.Before(eventsCollection[j].Timestamp) // sort by date ASC
	})

	log.Debugln("Events col", len(eventsCollection))

	return &eventsCollection
}
