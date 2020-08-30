package psqlstore //nolint:all

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/models"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SQLStore struct {
	conn *pgxpool.Pool
}

func NewPSQLStore(url string) (*SQLStore, error) {
	conn, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &SQLStore{conn}, nil
}

func (s *SQLStore) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *SQLStore) CreateEvent(ctx context.Context, e models.Event) error {
	if e.StartTime.After(e.EndTime) {
		return models.ErrTimeFrameInvalid
	}

	e.StartTime = e.StartTime.UTC().Round(time.Second)
	e.EndTime = e.EndTime.UTC().Round(time.Second)
	e.NotifyTime = e.NotifyTime.UTC().Round(time.Second)
	_, err := s.conn.Exec(ctx,
		"INSERT INTO events (title, description, start_time, end_time,notify_time,user_id) VALUES($1, $2, $3, $4, $5, $6) RETURNING event_id",
		e.Title, e.Description, e.StartTime, e.EndTime, e.NotifyTime, e.UserID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23P01" {
			return models.ErrTimeBusy
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *SQLStore) UpdateEvent(ctx context.Context, e models.Event) error {
	if e.StartTime.After(e.EndTime) {
		return models.ErrTimeFrameInvalid
	}
	_, err := s.conn.Exec(ctx,
		"UPDATE events SET title = $2, description = $3, start_time = $4,end_time = $5,notify_time = $6, user_id=$7 WHERE event_id = $1",
		e.ID, e.Title, e.Description, e.StartTime, e.EndTime, e.NotifyTime, e.UserID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23P01" {
			return models.ErrTimeBusy
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *SQLStore) GetEvent(ctx context.Context, id int) (models.Event, error) {
	var e models.Event
	err := s.conn.QueryRow(ctx, "select * from events where event_id=$1", id).Scan(
		&e.ID,
		&e.Title,
		&e.Description,
		&e.StartTime,
		&e.EndTime,
		&e.NotifyTime,
		&e.UserID,
	)
	if err != nil {
		return e, err
	}

	return e, nil
}

func (s *SQLStore) DeleteEvent(ctx context.Context, id int) error {
	tag, err := s.conn.Exec(ctx, "DELETE FROM events WHERE event_id = $1", id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return models.ErrEventNotFound
	}

	return nil
}

func (s *SQLStore) ListEvents(ctx context.Context, start time.Time, end time.Time) ([]models.Event, error) {
	if start.After(end) {
		return nil, models.ErrTimeFrameInvalid
	}
	rows, err := s.conn.Query(ctx,
		`SELECT *
		FROM events
		WHERE start_time > $1
	    AND start_time < $2
	    ORDER BY start_time
       `, start, end,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var eventList []models.Event
	for rows.Next() {
		var event models.Event
		err = rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.NotifyTime,
			&event.UserID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		eventList = append(eventList, event)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return eventList, nil
}
