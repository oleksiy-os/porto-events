package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/oleksiy-os/porto-events/internal/model"
	"github.com/oleksiy-os/porto-events/internal/store"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testDbPath = "tests_events_bolt.db"

var (
	logTestHook = NewTestLogger()
	repo        = &EventRepository{
		dbPath: testDbPath,
	}
)

func TestEventRepository_Add(t *testing.T) {
	if err := resetTestDb(repo); err != nil {
		t.Fatal("failed create test db|", err)
	}

	tests := []struct {
		name           string
		event          *model.Event
		wantAddOk      bool
		wantLogMessage bool
		logMessage     string
	}{
		{
			name:      "ok",
			wantAddOk: true,
			event: &model.Event{
				ID:          "event 1",
				Url:         "https://ev1.com",
				Title:       "Event 1",
				Description: "Ev 1 Description",
				Image:       "https://ev1.com/image.jpg",
				Place:       "Centre Square",
				LocationMap: "",
				Days:        "MO TU",
			},
			wantLogMessage: false,
		},
		{
			name:      "no add, already exists",
			wantAddOk: false,
			event: &model.Event{
				ID: "event 1",
			},
			wantLogMessage: true,
			logMessage:     "add event, already exists, will be skipped| event 1 ",
		},
		{
			name: "ok with id with spec chars",
			event: &model.Event{
				ID: `id with strange name | "with spec characters"`,
			},
			wantLogMessage: false,
			wantAddOk:      true,
		},
		{
			name: "already exist, skip adding",
			event: &model.Event{
				ID: `id with strange name | "with spec characters"`,
			},
			wantAddOk:      false,
			wantLogMessage: true,
			logMessage:     "add event, already exists, will be skipped| id with strange name | \"with spec characters\" ",
		},
	}
	for _, tt := range tests {
		logTestHook.Reset()
		t.Run(tt.name, func(t *testing.T) {
			repo.Add(tt.event)

			if tt.wantLogMessage {
				if logTestHook.LastEntry() == nil {
					t.Fatal("log is empty")
				}
				assert.Equal(t, tt.logMessage, logTestHook.LastEntry().Message)
				return
			}

			if len(logTestHook.Entries) > 0 {
				t.Error("called log", logTestHook.LastEntry().Message)
			}

			addedToMemory := false
			for _, e := range repo.events {
				if e.ID == tt.event.ID {
					addedToMemory = true
					break
				}
			}

			// check if event added to DB
			db, err := repo.openDb()
			defer closeDb(db)
			assert.Nil(t, err)
			//goland:noinspection GoUnhandledErrorResult
			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("Event"))
				if tt.wantAddOk {
					assert.NotNil(t, b.Get([]byte(tt.event.ID)), "not added to DB", tt.event.ID)
				} else {
					assert.Nil(t, b.Get([]byte(tt.event.ID)), "added to DB", tt.event.ID)
				}
				return nil
			})

			if tt.wantAddOk {
				assert.Truef(t, addedToMemory, "not added to repository")
			} else {
				assert.Falsef(t, addedToMemory, "should be NOT add to repository, but add")
			}
		})
	}
}

func TestEventRepository_Delete(t *testing.T) {
	repo.events = map[string]model.Event{
		"event 1": {
			ID:          "event 1",
			Url:         "https://ev1.com",
			Title:       "Event 1",
			Description: "Ev 1 Description",
			Image:       "https://ev1.com/image.jpg",
			Place:       "Centre Square",
			Category:    0,
		},
		"event 2": {
			ID:          "event 2",
			Url:         "https://ev2.com",
			Title:       "Event 2",
			Description: "Ev 2 Description",
			Image:       "https://ev1.com/image.jpg",
			Place:       "Centre Square",
			Days:        "MO TU",
			Category:    1,
		},
	}

	if err := resetTestDb(repo); err != nil {
		t.Fatal("failed create test db|", err)
	}

	tests := []struct {
		name       string
		wantOk     bool
		logMessage string
		id         string
	}{
		{
			name:   "ok",
			wantOk: true,
			id:     "event 1",
		},
		{
			name:   "not found",
			wantOk: false,
			id:     "event 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotOk := repo.GetById(tt.id)
			assert.Equal(t, tt.wantOk, gotOk, "not found id", tt.id)

			assert.Equal(t, tt.wantOk, repo.Delete(tt.id))

			_, gotOk = repo.GetById(tt.id)
			assert.False(t, gotOk, "not found eventById")
		})
	}
}

func TestEventRepository_ChangeCategory(t *testing.T) {
	repo.events = map[string]model.Event{
		"event 1": {
			ID:          "event 1",
			Url:         "https://ev1.com",
			Title:       "Event 1",
			Description: "Ev 1 Description",
			Image:       "https://ev1.com/image.jpg",
			Place:       "Centre Square",
			Category:    0,
		},
		"event 2": {
			ID:          "event 2",
			Url:         "https://ev2.com",
			Title:       "Event 2",
			Description: "Ev 2 Description",
			Image:       "https://ev1.com/image.jpg",
			Place:       "Centre Square",
			Days:        "MO TU",
			Category:    1,
		},
	}

	if err := resetTestDb(repo); err != nil {
		t.Fatal("failed create test db|", err)
	}

	tests := []struct {
		name       string
		wantOk     bool
		logMessage string
		args       store.ChangeCategoryData
	}{
		{
			name:   "change to publish",
			wantOk: true,
			args: store.ChangeCategoryData{
				Id:       "event 1",
				Category: 1,
			},
		},
		{
			name:   "change to new",
			wantOk: true,
			args: store.ChangeCategoryData{
				Id:       "event 2",
				Category: 0,
			},
		},
		{
			name:   "wrong category",
			wantOk: false,
			args: store.ChangeCategoryData{
				Id:       "event 2",
				Category: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evBeforeChange, _ := repo.GetById(tt.args.Id)
			assert.Equal(t, tt.wantOk, repo.ChangeCategory(tt.args))
			ev, ok := repo.GetById(tt.args.Id)
			assert.True(t, ok, "not found eventById")

			if tt.wantOk {
				assert.Equal(t, tt.args.Category, ev.Category)
			} else {
				assert.Equal(t, evBeforeChange.Category, ev.Category)
			}
		})
	}
}

func TestEventRepository_Save(t *testing.T) {
	if err := resetTestDb(repo); err != nil {
		t.Fatal("failed create test db|", err)
	}

	tests := []struct {
		name        string
		want        bool
		logMessage  string
		storedEvent *model.Event
		event       *model.Event
	}{
		{
			name: "ok",
			want: true,
			storedEvent: &model.Event{
				ID:          "event 1",
				Url:         "https://ev1-stored.com",
				Title:       "Event 1 - stored",
				Description: "Ev 1 Description - stored",
				Image:       "https://ev1.com/image-stored.jpg",
				Place:       "Centre Square-stored",
				Days:        "MO TU WE",
			},
			event: &model.Event{
				ID:          "event 1",
				Url:         "https://ev1.com",
				Title:       "Event 1",
				Description: "Ev 1 Description",
				Image:       "https://ev1.com/image.jpg",
				Place:       "Centre Square",
				Days:        "MO TU",
			},
		},
		{
			name:       "not found to save",
			want:       false,
			logMessage: "not found event to save|",
			event: &model.Event{
				ID:          "event 2",
				Url:         "https://ev2.com",
				Title:       "Event 2",
				Description: "Ev 2 Description",
				Image:       "https://ev1.com/image.jpg",
				Place:       "Centre Square",
				Days:        "MO TU",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want {
				repo.Add(tt.storedEvent) // store event to DB before check

				if assert.True(t, tt.want, repo.Save(tt.event), "save event", tt.event.ID) {
					ev, ok := repo.GetById(tt.event.ID)
					if !ok {
						t.Fatal("getById func returned false. LogMsg:", logTestHook.LastEntry().Message)
					}
					assert.Equal(t, tt.event.Title, ev.Title)
					assert.Equal(t, tt.event.Url, ev.Url)
					assert.Equal(t, tt.event.Description, ev.Description)
					assert.Equal(t, tt.event.Image, ev.Image)
					assert.Equal(t, tt.event.Place, ev.Place)
					assert.Equal(t, tt.event.Days, ev.Days)
				}
			}

			if !tt.want {
				assert.False(t, repo.Save(tt.event), "save event", tt.event.ID)
				assert.True(t, len(logTestHook.Entries) > 0, "log is empty")
				assert.Contains(t, logTestHook.LastEntry().Message, tt.logMessage)
			}
		})
	}
}

func Test_validateCategory(t *testing.T) {
	type args struct {
		category uint8
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
			args:    args{category: 0},
		},
		{
			name:    "ok",
			wantErr: false,
			args:    args{category: 1},
		},
		{
			name:    "not found",
			wantErr: true,
			args:    args{category: 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.EqualError(t, validateCategory(tt.args.category), "no such category")
			} else {
				assert.Nil(t, validateCategory(tt.args.category))
			}
		})
	}
}

func TestEventRepository_openDb(t *testing.T) {
	r := &EventRepository{}

	tests := []struct {
		name   string
		dbPath string
		want   bool
	}{
		{
			name:   "ok",
			dbPath: "events_bolt.db",
			want:   true,
		},
		{
			name:   "wrong db file path",
			dbPath: "wrong/path.db",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.dbPath = tt.dbPath
			_, err := r.openDb()
			assert.Equal(t, tt.want, err == nil, err)

		})
	}
}

func TestEventRepository_GetById(t *testing.T) {
	repo.events = map[string]model.Event{
		"event 1": {
			ID:          "event 1",
			Url:         "https://ev1-stored.com",
			Title:       "Event 1 - stored",
			Description: "Ev 1 Description - stored",
			Image:       "https://ev1.com/image-stored.jpg",
			Place:       "Centre Square-stored",
			Days:        "MO TU WE",
		},
		"event 2": {
			ID: "event 2",
		},
	}

	tests := []struct {
		name   string
		id     string
		wantEv model.Event
		wantOk bool
	}{
		{
			name:   "ok",
			id:     "event 1",
			wantOk: true,
			wantEv: model.Event{
				ID: "event 1",
			},
		},
		{
			name:   "not found",
			id:     "event x",
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := repo.GetById(tt.id)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantOk, got != nil)
			if tt.wantOk {
				assert.Equal(t, tt.wantEv.ID, got.ID)
			}
		})
	}
}

func resetTestDb(r *EventRepository) error {
	db, err := r.openDb()
	defer closeDb(db)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		//goland:noinspection GoUnhandledErrorResult
		tx.DeleteBucket([]byte("Event")) // clear events in test DB
		_, err = tx.CreateBucket([]byte("Event"))
		return err
	})

	return err
}

// NewTestLogger init test logger to catch all log messages during tests
func NewTestLogger() *test.Hook {
	// create test logger to catch log messages
	hook := test.NewGlobal()
	log.SetLevel(5)
	log.AddHook(hook)

	return hook
}

func TestEventRepository_GetCategoryPublish(t *testing.T) {
	type fields struct {
		events map[string]model.Event
	}
	tests := []struct {
		name   string
		fields fields
		want   []string // id category
	}{
		{
			name: "ok",
			fields: fields{
				events: map[string]model.Event{
					"1": {ID: "1", Category: 0},
					"2": {ID: "2", Category: 1},
					"3": {ID: "3", Category: 1},
					"4": {ID: "4", Category: 1},
				},
			},
			want: []string{"2", "3", "4"},
		},
		{
			name: "all publish",
			fields: fields{
				events: map[string]model.Event{
					"1": {ID: "1", Category: 1},
					"2": {ID: "2", Category: 1},
					"3": {ID: "3", Category: 1},
					"4": {ID: "4", Category: 1},
				},
			},
			want: []string{"1", "2", "3", "4"},
		},
		{
			name: "no publish",
			fields: fields{
				events: map[string]model.Event{
					"1": {ID: "1", Category: 0},
					"2": {ID: "2", Category: 0},
					"3": {ID: "3", Category: 0},
					"4": {ID: "4", Category: 0},
				},
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EventRepository{
				events: tt.fields.events,
			}

			gotEvs := r.GetCategoryPublish()

			assert.Equal(t, len(tt.want), len(*gotEvs))
			for _, want := range tt.want {
				ok := false
				for _, got := range *gotEvs {
					if got.ID == want {
						ok = true
						break
					}
				}
				assert.Truef(t, ok, "not equal got: %v and want %v", gotEvs, want)
			}
		})
	}
}
