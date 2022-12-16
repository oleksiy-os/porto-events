package store

import "github.com/oleksiy-os/porto-events/internal/model"

const (
	CategoryNew       = 0
	CategoryPublish   = 1
	CategoryPublished = 2
	CategoryBlocked   = 3
)

var CategoryId = map[uint8]uint8{
	CategoryNew:       CategoryNew,
	CategoryPublish:   CategoryPublish,
	CategoryPublished: CategoryPublished,
	CategoryBlocked:   CategoryBlocked,
}

type (
	ChangeCategoryData struct {
		Id       string
		Category uint8
	}

	EventRepository interface {
		// Get events from storage
		Get() *map[string]model.Event

		// GetById event from storage
		GetById(id string) (*model.Event, bool)

		// GetCategoryPublish list of events to publish
		GetCategoryPublish() *[]model.Event

		// Add events to storage.
		//
		// If failed log.Fatal will be called
		Add(*model.Event)

		// Save event to storage.
		Save(*model.Event) bool

		// Delete event
		Delete(string) bool

		// ChangeCategory event (new OR publish)
		ChangeCategory(data ChangeCategoryData) bool
	}
)
