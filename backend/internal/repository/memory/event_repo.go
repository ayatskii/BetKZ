package memory

import (
	"sync"

	"BetKZ/internal/domain"
)

type EventRepository struct {
	mu     sync.Mutex
	events map[int64]domain.Event
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		events: map[int64]domain.Event{
			1: {ID: 1, Name: "Real Madrid vs Barcelona", Odds: 1.75, Open: true},
			2: {ID: 2, Name: "UFC Main Event", Odds: 2.10, Open: true},
		},
	}
}

func (r *EventRepository) List() []domain.Event {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []domain.Event
	for _, e := range r.events {
		result = append(result, e)
	}
	return result
}

func (r *EventRepository) GetByID(id int64) (domain.Event, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	e, ok := r.events[id]
	return e, ok
}
