package models

import "errors"

var (
	ErrTimeBusy         = errors.New("timeframe busy")
	ErrTimeFrameInvalid = errors.New("invalid timeframe")
	ErrEventNotFound    = errors.New("event not found")
)
