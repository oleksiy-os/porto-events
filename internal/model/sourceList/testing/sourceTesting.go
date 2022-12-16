package testing

import (
	"encoding/json"
	m "github.com/oleksiy-os/porto-events/internal/model"
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
)

type (
	SourceTesting struct {
		pathToFile string
	}
)

func (s *SourceTesting) LoadEvents(_ *url.URL) []m.Event {
	var events []m.Event

	file, err := os.ReadFile(s.pathToFile)
	if err != nil {
		log.Error("wrong events file path|", err)
	}

	if err = json.Unmarshal(file, &events); err != nil {
		log.Error("json unmarshal|", err)
	}

	return events
}

func New() *SourceTesting {
	return &SourceTesting{
		pathToFile: "internal/model/sourceList/testing/test_events_list.json",
	}
}
