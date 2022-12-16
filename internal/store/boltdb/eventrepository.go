package boltdb

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/oleksiy-os/porto-events/internal/model"
	"github.com/oleksiy-os/porto-events/internal/store"
	log "github.com/sirupsen/logrus"
)

const dbPath = "internal/store/boltdb/events_bolt.db"

type (
	EventRepository struct {
		dbPath string
		events map[string]model.Event
	}
)

func (r *EventRepository) Get() *map[string]model.Event {
	if r.events != nil {
		return &r.events
	}

	r.events = make(map[string]model.Event)

	db, err := r.openDb()
	defer closeDb(db)
	if err != nil {
		log.Error("db access|", err)
		return &r.events
	}

	if err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Event"))
		if b == nil {
			b, err = createBucket(db)
			if err != nil {
				return err
			}
		}

		var val model.Event
		if err = b.ForEach(func(k, v []byte) error {
			if err = json.Unmarshal(v, &val); err != nil {
				log.Error("decode bolt|", err)
			}
			r.events[string(k)] = val
			return nil
		}); err != nil {
			return err
		}
		return err
	}); err != nil {
		log.Error("get events|", err)
	}

	return &r.events
}

func (r *EventRepository) GetById(id string) (*model.Event, bool) {
	if r.events == nil {
		r.Get()
	}

	if r.events[id].ID == "" {
		return nil, false
	}
	ev := r.events[id]
	return &ev, true
}

func (r *EventRepository) GetCategoryPublish() *[]model.Event {
	var evsPublish []model.Event
	for _, e := range *r.Get() {
		if e.Category == store.CategoryPublish {
			evsPublish = append(evsPublish, e)
		}
	}

	return &evsPublish
}

func (r *EventRepository) Add(event *model.Event) {
	if _, isExist := r.GetById(event.ID); isExist {
		log.Debugln("add event, already exists, will be skipped|", event.ID, event.Title)
		return
	}

	if event.ID == "" {
		event.ID = event.Title
	}

	db, err := r.openDb()
	defer closeDb(db)
	if err != nil {
		log.Error("add event|", err)
		return
	}

	dbPut(db, event)

	r.events[event.ID] = *event
}

func (r *EventRepository) Save(event *model.Event) bool {
	if _, ok := r.GetById(event.ID); !ok {
		log.Error("not found event to save|", event)
		return false
	}

	db, err := r.openDb()
	defer closeDb(db)
	if err != nil {
		log.Error("save event|", err)
		return false
	}

	if ok := dbPut(db, event); !ok {
		return false
	}

	r.events[event.ID] = *event

	return true
}

func (r *EventRepository) Delete(id string) bool {
	if _, ok := r.GetById(id); !ok {
		log.Error("not found event to delete|", id)
		return false
	}

	db, err := r.openDb()
	defer closeDb(db)
	if err != nil {
		log.Error("delete event|", err)
		return false
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		b = tx.Bucket([]byte("Event"))

		return b.Delete([]byte(id))
	}); err != nil {
		log.Error("delete event|", err)
		return false
	}
	delete(r.events, id)

	return true
}

func (r *EventRepository) ChangeCategory(d store.ChangeCategoryData) bool {
	d.Id = model.StripAllHtml.Sanitize(d.Id)
	db, err := r.openDb()
	defer closeDb(db)
	if err != nil {
		log.Error("change category|", err)
		return false
	}

	if err = validateCategory(d.Category); err != nil {
		log.Error("save events|", err)
		return false
	}

	event := r.events[d.Id]
	event.Category = d.Category

	if ok := dbPut(db, &event); !ok {
		return false
	}

	r.events[d.Id] = event
	return true
}

func dbPut(db *bolt.DB, event *model.Event) bool {
	if err := db.Update(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		b = tx.Bucket([]byte("Event"))

		if event.ID == "" {
			return errors.New("not found. id: " + event.Title)
		}

		evJson, err := json.Marshal(event)
		if err != nil {
			return err
		}

		return b.Put([]byte(event.ID), evJson)
	}); err != nil {
		log.Error("save events|", err)
		return false
	}

	return true
}

func createBucket(db *bolt.DB) (*bolt.Bucket, error) {
	var b *bolt.Bucket

	err := db.Update(func(tx *bolt.Tx) error {
		var err error
		b, err = tx.CreateBucketIfNotExists([]byte("Event"))
		return err
	})

	return b, err
}

func validateCategory(category uint8) error {
	for _, c := range store.CategoryId {
		if category == c {
			return nil
		}
	}
	return errors.New("no such category")
}

func closeDb(db *bolt.DB) {
	if err := db.Close(); err != nil {
		log.Fatal("close db|", err)
	}
}

func (r *EventRepository) openDb() (*bolt.DB, error) {
	if r.dbPath == "" {
		r.dbPath = dbPath
	}

	return bolt.Open(r.dbPath, 0600, nil)
}
