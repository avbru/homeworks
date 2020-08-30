package models

import "time"

type Event struct {
	ID          int
	Title       string
	StartTime   time.Time
	EndTime     time.Time
	Description string
	UserID      string
	NotifyTime  time.Time
}
