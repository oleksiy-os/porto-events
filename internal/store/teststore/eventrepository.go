package teststore

import (
	"github.com/oleksiy-os/porto-events/internal/model"
	"github.com/oleksiy-os/porto-events/internal/store"
	log "github.com/sirupsen/logrus"
)

type (
	TestEventRepository struct {
		dbPath string
		events map[string]model.Event
	}
)

func (r *TestEventRepository) GetCategoryPublish() *[]model.Event {
	//TODO implement me
	log.Error("not implemented teststore/GetCategoryPublish")

	var events *[]model.Event
	return events
}

func (r *TestEventRepository) Get() *map[string]model.Event {
	if r.events != nil {
		return &r.events
	}

	r.events = make(map[string]model.Event)

	return &r.events
}

func (r *TestEventRepository) GetById(id string) (*model.Event, bool) {
	if r.events == nil {
		r.Get()
	}

	if r.events[id].ID == "" {
		return nil, false
	}
	ev := r.events[id]
	return &ev, true
}

func (r *TestEventRepository) Add(event *model.Event) {
	if _, isExist := r.GetById(event.ID); isExist {
		log.Debugln("add event, already exists, will be skipped|", event.ID, event.Title)
		return
	}

	if event.ID == "" {
		event.ID = event.Title
	}

	r.events[event.ID] = *event
}

func (r *TestEventRepository) Save(event *model.Event) bool {
	if _, ok := r.GetById(event.ID); !ok {
		log.Error("not found event to save|", event)
		return false
	}

	r.events[event.ID] = *event

	return true
}

func (r *TestEventRepository) Delete(id string) bool {
	if _, ok := r.GetById(id); !ok {
		log.Error("not found event to delete|", id)
		return false
	}

	delete(r.events, id)

	return true
}

func (r *TestEventRepository) ChangeCategory(d store.ChangeCategoryData) bool {
	d.Id = model.StripAllHtml.Sanitize(d.Id)

	event := r.events[d.Id]
	event.Category = d.Category

	r.events[d.Id] = event
	return true
}
