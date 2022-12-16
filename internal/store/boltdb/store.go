package boltdb

import (
	"github.com/oleksiy-os/porto-events/internal/store"
)

type Store struct {
	eventRepository *EventRepository
}

func (s *Store) Event() store.EventRepository {
	return s.eventRepository
}

func New() *Store {
	return &Store{
		eventRepository: &EventRepository{},
	}
}
