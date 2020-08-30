package storage

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/models"
	"github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/storages/psqlstore"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	url := "postgres://calendar:derparol@devdb:5432/calendar?sslmode=disable"
	if err := Migrate(url); err != nil {
		log.Fatalf("%s", err)
	}
	s, err := psqlstore.NewPSQLStore(url)
	require.NoError(t, err)

	baseTime := time.Now().UTC().Round(time.Second)
	event := models.Event{
		Title:       "title",
		StartTime:   baseTime,
		EndTime:     baseTime.Add(time.Minute),
		Description: "description",
		UserID:      "userid",
		NotifyTime:  baseTime.Add(-time.Minute),
	}

	err = s.CreateEvent(context.Background(), event) // 1 event IDs 1
	require.NoError(t, err)

	// Create TimeFrame Invalid
	err = s.CreateEvent(context.Background(), models.Event{StartTime: baseTime.Add(time.Hour), EndTime: baseTime})
	require.Equal(t, models.ErrTimeFrameInvalid, err)

	// Create Timeframe not available
	err = s.CreateEvent(context.Background(), event)
	require.Equal(t, models.ErrTimeBusy, err)

	// Create Timeframe  available
	event2 := event
	event2.StartTime = event.StartTime.Add(time.Hour)
	event2.EndTime = event.EndTime.Add(time.Hour)
	err = s.CreateEvent(context.Background(), event2) // 2 events IDs 1,2
	require.NoError(t, err)

	e, err := s.GetEvent(context.Background(), 1)
	event.ID = 1
	require.Equal(t, event, e)

	// Create TimeFrame Invalid
	err = s.UpdateEvent(context.Background(), models.Event{StartTime: baseTime.Add(time.Hour), EndTime: baseTime})
	require.Equal(t, models.ErrTimeFrameInvalid, err)

	// Update Timeframe not available
	e.EndTime = e.EndTime.Add(time.Hour)
	err = s.UpdateEvent(context.Background(), e)
	require.Equal(t, models.ErrTimeBusy, err)

	e.StartTime = event.StartTime.Add(time.Hour * 2)
	e.EndTime = event.EndTime.Add(time.Hour * 2)
	err = s.UpdateEvent(context.Background(), e)
	require.NoError(t, err)

	// ListEvents Invalid Timeframe
	events, err := s.ListEvents(context.Background(), baseTime.Add(time.Hour), baseTime)
	require.Equal(t, models.ErrTimeFrameInvalid, err)

	events, err = s.ListEvents(context.Background(), baseTime.Add(-time.Hour), baseTime.Add(time.Hour*10))
	require.Equal(t, 2, len(events))

	// Delete non existing
	err = s.DeleteEvent(context.Background(), 30)
	require.Error(t, err)

	// Delete
	err = s.DeleteEvent(context.Background(), 1)
	require.NoError(t, err)

	events, err = s.ListEvents(context.Background(), baseTime.Add(-time.Hour), baseTime.Add(time.Hour*10))
	require.Equal(t, 1, len(events))
}

func Migrate(url string) error {
	m, err := migrate.New(
		"file:///app/migrations",
		url)
	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}
	return nil
}
