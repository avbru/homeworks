package memstore

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/models"

	"github.com/stretchr/testify/require"
)

func TestMemStore(t *testing.T) {
	s, _ := NewEventStore()
	baseTime := time.Now().UTC().Round(time.Second)
	event := models.Event{
		Title:       "title",
		StartTime:   baseTime,
		EndTime:     baseTime.Add(time.Minute),
		Description: "description",
		UserID:      "userid",
		NotifyTime:  baseTime.Add(-time.Minute),
	}

	err := s.CreateEvent(context.Background(), event) // 1 event IDs 1
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
	err = s.DeleteEvent(context.Background(), 3)
	require.Error(t, err)

	// Delete
	err = s.DeleteEvent(context.Background(), 1)
	require.NoError(t, err)

	events, err = s.ListEvents(context.Background(), baseTime.Add(-time.Hour), baseTime.Add(time.Hour*10))
	require.Equal(t, 1, len(events))
}

func TestMemStoreConcurrent(t *testing.T) {
	s, _ := NewEventStore()
	wg := &sync.WaitGroup{}
	wg.Add(4)

	done := make(chan struct{})

	baseTime := time.Now()
	// Create
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			if i == 300 {
				close(done)
			}

			_ = s.CreateEvent(context.Background(), models.Event{
				ID:          i,
				Title:       strconv.Itoa(i),
				StartTime:   baseTime.Add(time.Minute * time.Duration(i)),
				EndTime:     baseTime.Add(time.Minute*time.Duration(i) + 1),
				Description: strconv.Itoa(i),
				UserID:      "userid",
				NotifyTime:  baseTime.Add(-time.Minute * 1),
			})
		}
	}()

	<-done

	// Delete
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i = i + 3 {
			_ = s.DeleteEvent(context.Background(), i)
		}
	}()

	// Get, Update
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			e, _ := s.GetEvent(context.Background(), i)
			e.Description = "updated description"
			_ = s.UpdateEvent(context.Background(), e)
		}
	}()

	// List
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_, _ = s.ListEvents(context.Background(), baseTime, baseTime.Add(time.Minute*1000))
		}
	}()

	wg.Wait()
}

func TestTimeAvailable(t *testing.T) {
	baseTime := time.Now()
	event := models.Event{
		StartTime: baseTime,
		EndTime:   baseTime.Add(time.Minute * 2),
	}
	s, _ := NewEventStore()
	err := s.CreateEvent(context.Background(), event)
	require.NoError(t, err)

	// Exactly same time
	require.False(t, s.isTimeAvailable(event))

	// Start time is in busy timeframe
	event.StartTime = event.StartTime.Add(time.Minute)
	require.False(t, s.isTimeAvailable(event))

	// End time is in busy timeframe
	event.StartTime = baseTime.Add(-time.Minute)
	event.EndTime = baseTime.Add(time.Minute)
	require.False(t, s.isTimeAvailable(event))

	// New event overlaps event from both sides
	event.StartTime = baseTime.Add(-time.Minute)
	event.EndTime = baseTime.Add(time.Minute * 3)
	require.False(t, s.isTimeAvailable(event))

	// No overlap before existing timeframe
	event.StartTime = baseTime.Add(-time.Minute * 2)
	event.EndTime = baseTime.Add(-time.Minute)
	require.True(t, s.isTimeAvailable(event))

	// No overlap after existing timeframe
	event.StartTime = baseTime.Add(time.Minute * 3)
	event.EndTime = baseTime.Add(time.Minute * 4)
	require.True(t, s.isTimeAvailable(event))
}
