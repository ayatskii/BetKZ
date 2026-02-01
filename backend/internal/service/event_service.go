package service

import "BetKZ/internal/domain"

type EventService struct {
	repo interface {
		List() []domain.Event
	}
}

func NewEventService(r interface {
	List() []domain.Event
}) *EventService {
	return &EventService{repo: r}
}

func (s *EventService) ListEvents() []domain.Event {
	return s.repo.List()
}
