package model

import (
	"time"
)

// TimeOfDay represents a time of day without a date component.
// Used for nested JSON objects like {"hours": 13, "minutes": 36, "seconds": 43}
type TimeOfDay struct {
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
	Seconds int `json:"seconds"`
}

func (t TimeOfDay) CombineWith(date time.Time) time.Time {
	return time.Date(
		date.Year(), date.Month(), date.Day(),
		t.Hours, t.Minutes, t.Seconds, 0,
		date.Location(),
	)
}
