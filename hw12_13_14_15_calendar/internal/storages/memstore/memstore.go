package memstore

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/models"
)

type EventList []models.Event

type EventStore struct {
	sync.RWMutex
	maxID  int
	events map[int]models.Event
}

func NewEventStore() (*EventStore, error) {
	return &EventStore{
		events: make(map[int]models.Event),
	}, nil
}

func (s *EventStore) CreateEvent(ctx context.Context, event models.Event) error {
	s.Lock()
	defer s.Unlock()

	if event.StartTime.After(event.EndTime) {
		return models.ErrTimeFrameInvalid
	}

	if !s.isTimeAvailable(event) {
		return models.ErrTimeBusy
	}

	s.maxID++
	event.ID = s.maxID
	s.events[event.ID] = event

	return nil
}

func (s *EventStore) UpdateEvent(ctx context.Context, event models.Event) error {
	if event.StartTime.After(event.EndTime) {
		return models.ErrTimeFrameInvalid
	}

	s.Lock()
	defer s.Unlock()
	if _, ok := s.events[event.ID]; !ok {
		return models.ErrEventNotFound
	}

	if !s.isTimeAvailable(event) {
		return models.ErrTimeBusy
	}

	s.events[event.ID] = event

	return nil
}

func (s *EventStore) GetEvent(ctx context.Context, id int) (models.Event, error) {
	s.Lock()
	defer s.Unlock()
	e, ok := s.events[id]
	if !ok {
		return models.Event{}, nil
	}

	return e, nil
}

func (s *EventStore) DeleteEvent(ctx context.Context, id int) error {
	s.Lock()
	defer s.Unlock()
	_, ok := s.events[id]
	if !ok {
		return fmt.Errorf("delete event error: %w", models.ErrEventNotFound)
	}
	delete(s.events, id)

	return nil
}

func (s *EventStore) ListEvents(ctx context.Context, start time.Time, end time.Time) ([]models.Event, error) {
	if start.After(end) {
		return nil, models.ErrTimeFrameInvalid
	}

	s.RLock()

	var events EventList
	events = make([]models.Event, 0)
	for _, e := range s.events {
		if e.StartTime.After(start) && e.StartTime.Before(end) {
			events = append(events, e)
		}
	}

	s.RUnlock()
	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime.Unix() < events[j].StartTime.Unix()
	})

	return events, nil
}

// isTimeAvailable mimics PSQL overlapping timeframes constraint.
func (s *EventStore) isTimeAvailable(event models.Event) bool {
	for _, e := range s.events {
		if e.ID == event.ID {
			continue
		}
		if !(event.EndTime.Before(e.StartTime) || event.StartTime.After(e.EndTime)) {
			return false
		}
	}

	return true
}
