package model

import (
	"github.com/BurntSushi/toml"
	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type (
	Source struct {
		Name string `toml:"name"`
		Url  string `toml:"url"`
	}

	Event struct {
		ID          string
		Url         string
		Title       string
		Description string
		Image       string
		Place       string
		Location    string
		LocationMap string
		DateText    string    // Example: "May 12th, 2022 - December 31st, 2022"
		Days        string    // working days. Ex.: "mon, tue, wed, thu, fri, sat, sun"
		Time        string    // "10:00 - 18:00"
		Timestamp   time.Time // for events sort, if event has date range, here will be current date
		Category    uint8     // 0: New event; 1: publish
	}
)

var StripAllHtml = bluemonday.StrictPolicy()

func GetSources(confPath string) ([]Source, error) {
	{
		var sources map[string][]Source
		if _, err := toml.DecodeFile(confPath, &sources); err != nil {
			return nil, err
		}
		return sources["source"], nil
	}
}

func LoadContent(res *http.Response) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(res.Body)
}

//goland:noinspection GoUnusedExportedFunction
func GetBaseUrl(sourceUrl string) (string, bool) {
	u, err := url.Parse(sourceUrl)
	if err != nil {
		log.Errorln("error GetBaseUrl", err)
		return "", false
	}

	return u.Scheme + "://" + u.Host, true
}
