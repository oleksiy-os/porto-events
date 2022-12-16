package teststore

import (
	"github.com/oleksiy-os/porto-events/internal/store"
)

type Store struct {
	eventRepository *TestEventRepository
}

func (s *Store) Event() store.EventRepository {
	return s.eventRepository
}

func New() *Store {
	return &Store{
		eventRepository: &TestEventRepository{},
	}
}
