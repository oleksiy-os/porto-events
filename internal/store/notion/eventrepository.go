package notion

// TODO: Notion integration not ready to use. Under development

import (
	notion "github.com/jomei/notionapi"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type (
	NotiRepository struct {
		client *notion.Client
		config Notion
	}
)

func (r *NotiRepository) Save() {}
func (r *NotiRepository) Get() {
	query, err := r.client.Database.Query(context.Background(), notion.DatabaseID(r.config.PageEventsId), nil)
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, page := range query.Results {
		for _, prop := range page.Properties {
			//if prop.GetType() == "title" {
			//	prop.
			//break
			//}

			log.Info("props", prop)
		}
	}
	//content := block.Results[0].(*notion.ParagraphBlock).Paragraph.RichText[0].Text.Content

	log.Debugf("content: %+v", query)
}

// CheckEvery how often check source for new events, return value in hours
func (r *NotiRepository) CheckEvery() uint8 {
	var timer uint8 = 0

	block, err := r.client.Database.Query(context.Background(), notion.DatabaseID(r.config.PageConfigId), nil)
	if err != nil {
		log.Errorln("get conf error", err)
	}

	for _, result := range block.Results {
		title := result.Properties["Name"].(*notion.TitleProperty).Title
		if len(title) == 1 && title[0].PlainText == "timer" {
			timer = uint8(result.Properties["Number"].(*notion.NumberProperty).Number)
		}
	}

	if timer != 0 {
		log.Debugln("timer value found in notion config page and changed", timer)
		r.config.Timer = timer
	}

	return r.config.Timer
}

func (r *NotiRepository) GetPage() {
	page, err := r.client.Page.Get(context.Background(), notion.PageID(r.config.PageEventsId))
	if err != nil {
		log.Errorln(err)
	}
	var data string
	if err = page.Properties.UnmarshalJSON([]byte(data)); err != nil {
		log.Errorln(err)
	}
	log.Printf("data %+v", data)
	log.Printf("page %+v", page)
}
