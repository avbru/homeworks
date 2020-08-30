package storages

import (
	"context"
	"time"

	"github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/models"
)

type EventStore interface {
	CreateEvent(ctx context.Context, event models.Event) error
	UpdateEvent(ctx context.Context, event models.Event) error
	GetEvent(ctx context.Context, id int) (models.Event, error)
	DeleteEvent(ctx context.Context, id int) error
	ListEvents(ctx context.Context, start time.Time, end time.Time) ([]models.Event, error)
}
