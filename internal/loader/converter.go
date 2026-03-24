package loader

import (
	"fmt"
	"time"
)

var timeFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.000Z",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05",
	"2006-01-02",
}

func ParseTime(s string) (time.Time, error) {
	s = trimSpaces(s)
	if s == "" {
		return time.Time{}, nil
	}
	for _, layout := range timeFormats {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse time %q", s)
}

func ParseTimePtr(s string) (*time.Time, error) {
	s = trimSpaces(s)
	if s == "" {
		return nil, nil
	}
	t, err := ParseTime(s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func ParseStringPtr(s string) *string {
	s = trimSpaces(s)
	if s == "" {
		return nil
	}
	return &s
}

func ValidateRequired(value string, fieldName string) error {
	if trimSpaces(value) == "" {
		return fmt.Errorf("%s: required field is empty", fieldName)
	}
	return nil
}

func (r *LoadResult) String() string {
	return fmt.Sprintf("%s: total=%d inserted=%d skipped=%d errors=%d",
		r.Entity, r.Total, r.Inserted, r.Skipped, len(r.Errors))
}
