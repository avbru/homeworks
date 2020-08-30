package app

import (
	"context"

	"github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/models"
	"github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/storages"
)

type App struct {
	storage storages.EventStore
}

func New(storage storages.EventStore) *App {
	return &App{
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id int, title string) error {
	return a.storage.CreateEvent(ctx, models.Event{ID: id, Title: title})
}
