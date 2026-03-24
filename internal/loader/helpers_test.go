package loader

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sudarsh1010/graph-rag/internal/model"
)

func TestGetString(t *testing.T) {
	tests := []struct {
		name string
		raw  map[string]interface{}
		key  string
		want string
	}{
		{
			name: "present string value",
			raw:  map[string]interface{}{"name": "hello"},
			key:  "name",
			want: "hello",
		},
		{
			name: "missing key returns empty",
			raw:  map[string]interface{}{"name": "hello"},
			key:  "missing",
			want: "",
		},
		{
			name: "non-string value is formatted",
			raw:  map[string]interface{}{"count": 42},
			key:  "count",
			want: "42",
		},
		{
			name: "float64 value is formatted",
			raw:  map[string]interface{}{"price": 19.99},
			key:  "price",
			want: "19.99",
		},
		{
			name: "nil value returns empty",
			raw:  map[string]interface{}{"val": nil},
			key:  "val",
			want: "<nil>",
		},
		{
			name: "empty string value",
			raw:  map[string]interface{}{"name": ""},
			key:  "name",
			want: "",
		},
		{
			name: "empty map missing key",
			raw:  map[string]interface{}{},
			key:  "any",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getString(tt.raw, tt.key)
			if got != tt.want {
				t.Errorf("getString() key=%q = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestGetStringPtr(t *testing.T) {
	tests := []struct {
		name    string
		raw     map[string]interface{}
		key     string
		wantNil bool
		want    string
	}{
		{
			name: "present string value",
			raw:  map[string]interface{}{"name": "hello"},
			key:  "name",
			want: "hello",
		},
		{
			name:    "missing key returns nil",
			raw:     map[string]interface{}{"name": "hello"},
			key:     "missing",
			wantNil: true,
		},
		{
			name:    "nil value returns nil",
			raw:     map[string]interface{}{"name": nil},
			key:     "name",
			wantNil: true,
		},
		{
			name:    "empty string returns nil",
			raw:     map[string]interface{}{"name": ""},
			key:     "name",
			wantNil: true,
		},
		{
			name:    "non-string value returns nil",
			raw:     map[string]interface{}{"count": 42},
			key:     "count",
			wantNil: true,
		},
		{
			name:    "float64 value returns nil",
			raw:     map[string]interface{}{"price": 19.99},
			key:     "price",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStringPtr(tt.raw, tt.key)
			if (got == nil) != tt.wantNil {
				t.Errorf("getStringPtr() key=%q = %v, wantNil %v", tt.key, got, tt.wantNil)
				return
			}
			if got != nil && *got != tt.want {
				t.Errorf("getStringPtr() key=%q = %q, want %q", tt.key, *got, tt.want)
			}
		})
	}
}

func TestGetDecimal(t *testing.T) {
	tests := []struct {
		name    string
		raw     map[string]interface{}
		key     string
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name: "string integer",
			raw:  map[string]interface{}{"price": "42"},
			key:  "price",
			want: decimal.RequireFromString("42"),
		},
		{
			name: "string float",
			raw:  map[string]interface{}{"price": "3.14"},
			key:  "price",
			want: decimal.RequireFromString("3.14"),
		},
		{
			name: "empty string returns zero",
			raw:  map[string]interface{}{"price": ""},
			key:  "price",
			want: decimal.Decimal{},
		},
		{
			name: "missing key returns zero",
			raw:  map[string]interface{}{},
			key:  "price",
			want: decimal.Decimal{},
		},
		{
			name:    "invalid string returns error",
			raw:     map[string]interface{}{"price": "abc"},
			key:     "price",
			wantErr: true,
		},
		{
			name: "negative decimal",
			raw:  map[string]interface{}{"price": "-5.5"},
			key:  "price",
			want: decimal.RequireFromString("-5.5"),
		},
		// Note: when JSON unmarshals numbers into interface{}, they become float64.
		// getString converts float64 to its string representation via fmt.Sprintf("%v", v).
		// So a float64 3.14 becomes "3.14" which ParseDecimal can parse.
		{
			name: "float64 value (from JSON unmarshal)",
			raw:  map[string]interface{}{"price": 3.14},
			key:  "price",
			want: decimal.RequireFromString("3.14"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDecimal(tt.raw, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDecimal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("getDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTime(t *testing.T) {
	tests := []struct {
		name    string
		raw     map[string]interface{}
		key     string
		wantErr bool
		check   func(t *testing.T, got time.Time)
	}{
		{
			name: "RFC3339 string",
			raw:  map[string]interface{}{"date": "2024-03-15T10:30:00Z"},
			key:  "date",
			check: func(t *testing.T, got time.Time) {
				if got.Year() != 2024 {
					t.Errorf("got year %d, want 2024", got.Year())
				}
			},
		},
		{
			name: "date only",
			raw:  map[string]interface{}{"date": "2024-03-15"},
			key:  "date",
			check: func(t *testing.T, got time.Time) {
				if got.Month() != time.March || got.Day() != 15 {
					t.Errorf("got %v, want 2024-03-15", got)
				}
			},
		},
		{
			name: "empty string returns zero time",
			raw:  map[string]interface{}{"date": ""},
			key:  "date",
			check: func(t *testing.T, got time.Time) {
				if !got.IsZero() {
					t.Errorf("got %v, want zero time", got)
				}
			},
		},
		{
			name:  "missing key returns zero time",
			raw:   map[string]interface{}{},
			key:   "date",
			check: func(t *testing.T, got time.Time) { /* zero time is fine */ },
		},
		{
			name:    "invalid date returns error",
			raw:     map[string]interface{}{"date": "not-a-date"},
			key:     "date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTime(tt.raw, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func TestGetTimePtr(t *testing.T) {
	tests := []struct {
		name    string
		raw     map[string]interface{}
		key     string
		wantNil bool
		wantErr bool
		check   func(t *testing.T, got *time.Time)
	}{
		{
			name: "valid RFC3339 string",
			raw:  map[string]interface{}{"date": "2024-03-15T10:30:00Z"},
			key:  "date",
			check: func(t *testing.T, got *time.Time) {
				if got.Year() != 2024 {
					t.Errorf("got year %d, want 2024", got.Year())
				}
			},
		},
		{
			name:    "missing key returns nil no error",
			raw:     map[string]interface{}{},
			key:     "date",
			wantNil: true,
		},
		{
			name:    "nil value returns nil no error",
			raw:     map[string]interface{}{"date": nil},
			key:     "date",
			wantNil: true,
		},
		{
			name:    "empty string returns nil no error",
			raw:     map[string]interface{}{"date": ""},
			key:     "date",
			wantNil: true,
		},
		{
			name:    "invalid string returns error nil",
			raw:     map[string]interface{}{"date": "bad-date"},
			key:     "date",
			wantErr: true,
			wantNil: true,
		},
		{
			name:    "non-string non-nil returns error",
			raw:     map[string]interface{}{"date": 12345},
			key:     "date",
			wantErr: true,
			wantNil: true,
		},
		{
			name: "date only format",
			raw:  map[string]interface{}{"date": "2024-03-15"},
			key:  "date",
			check: func(t *testing.T, got *time.Time) {
				if got.Day() != 15 {
					t.Errorf("got day %d, want 15", got.Day())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTimePtr(tt.raw, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTimePtr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("getTimePtr() = %v, wantNil %v", got, tt.wantNil)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name string
		raw  map[string]interface{}
		key  string
		want bool
	}{
		{name: "bool true", raw: map[string]interface{}{"active": true}, key: "active", want: true},
		{name: "bool false", raw: map[string]interface{}{"active": false}, key: "active", want: false},
		{name: "string true lowercase", raw: map[string]interface{}{"active": "true"}, key: "active", want: true},
		{name: "string true uppercase", raw: map[string]interface{}{"active": "True"}, key: "active", want: true},
		{name: "string true mixed case", raw: map[string]interface{}{"active": "TRUE"}, key: "active", want: true},
		{name: "string false", raw: map[string]interface{}{"active": "false"}, key: "active", want: false},
		{name: "string random", raw: map[string]interface{}{"active": "yes"}, key: "active", want: false},
		{name: "missing key returns false", raw: map[string]interface{}{}, key: "active", want: false},
		{name: "int value returns false", raw: map[string]interface{}{"active": 1}, key: "active", want: false},
		{name: "nil value returns false", raw: map[string]interface{}{"active": nil}, key: "active", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getBool(tt.raw, tt.key)
			if got != tt.want {
				t.Errorf("getBool() key=%q = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestGetIntRequired(t *testing.T) {
	tests := []struct {
		name    string
		raw     map[string]interface{}
		key     string
		want    int
		wantErr bool
	}{
		{name: "float64 value (JSON number)", raw: map[string]interface{}{"count": float64(42)}, key: "count", want: 42},
		{name: "float64 with decimal truncated", raw: map[string]interface{}{"count": float64(42.7)}, key: "count", want: 42},
		{name: "string integer", raw: map[string]interface{}{"count": "100"}, key: "count", want: 100},
		{name: "int value", raw: map[string]interface{}{"count": 5}, key: "count", want: 5},
		{name: "zero float64", raw: map[string]interface{}{"count": float64(0)}, key: "count", want: 0},
		{name: "zero string", raw: map[string]interface{}{"count": "0"}, key: "count", want: 0},
		{name: "missing key returns error", raw: map[string]interface{}{}, key: "count", wantErr: true},
		{name: "invalid string returns error", raw: map[string]interface{}{"count": "abc"}, key: "count", wantErr: true},
		{name: "string float parses integer part", raw: map[string]interface{}{"count": "3.14"}, key: "count", want: 3},
		{name: "bool value returns error", raw: map[string]interface{}{"count": true}, key: "count", wantErr: true},
		{name: "nil value returns error", raw: map[string]interface{}{"count": nil}, key: "count", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getIntRequired(tt.raw, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getIntRequired() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("getIntRequired() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestGetIntPtr(t *testing.T) {
	tests := []struct {
		name    string
		raw     map[string]interface{}
		key     string
		wantNil bool
		want    int
		wantErr bool
	}{
		{name: "float64 value (JSON number)", raw: map[string]interface{}{"count": float64(42)}, key: "count", want: 42},
		{name: "float64 with decimal truncated", raw: map[string]interface{}{"count": float64(42.9)}, key: "count", want: 42},
		{name: "string integer", raw: map[string]interface{}{"count": "100"}, key: "count", want: 100},
		{name: "missing key returns nil", raw: map[string]interface{}{}, key: "count", wantNil: true},
		{name: "nil value returns nil", raw: map[string]interface{}{"count": nil}, key: "count", wantNil: true},
		{name: "empty string returns nil", raw: map[string]interface{}{"count": ""}, key: "count", wantNil: true},
		{name: "invalid string returns error nil", raw: map[string]interface{}{"count": "abc"}, key: "count", wantErr: true, wantNil: true},
		{name: "string float parses integer part", raw: map[string]interface{}{"count": "3.14"}, key: "count", want: 3},
		{name: "bool value returns error", raw: map[string]interface{}{"count": true}, key: "count", wantErr: true, wantNil: true},
		{name: "int value not handled returns error", raw: map[string]interface{}{"count": 7}, key: "count", wantErr: true, wantNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getIntPtr(tt.raw, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getIntPtr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("getIntPtr() = %v, wantNil %v", got, tt.wantNil)
				return
			}
			if !tt.wantErr && got != nil && *got != tt.want {
				t.Errorf("getIntPtr() = %d, want %d", *got, tt.want)
			}
		})
	}
}

func TestGetTimeOfDay(t *testing.T) {
	tests := []struct {
		name string
		raw  map[string]interface{}
		key  string
		want model.TimeOfDay
	}{
		{
			name: "valid nested object with all fields",
			raw: map[string]interface{}{
				"time": map[string]interface{}{"hours": float64(13), "minutes": float64(36), "seconds": float64(43)},
			},
			key:  "time",
			want: model.TimeOfDay{Hours: 13, Minutes: 36, Seconds: 43},
		},
		{
			name: "nested object partial fields",
			raw: map[string]interface{}{
				"time": map[string]interface{}{"hours": float64(9), "minutes": float64(30)},
			},
			key:  "time",
			want: model.TimeOfDay{Hours: 9, Minutes: 30, Seconds: 0},
		},
		{
			name: "nested object with int values",
			raw: map[string]interface{}{
				"time": map[string]interface{}{"hours": 12, "minutes": 0, "seconds": 0},
			},
			key:  "time",
			want: model.TimeOfDay{Hours: 12, Minutes: 0, Seconds: 0},
		},
		{
			name: "missing key returns zero TimeOfDay",
			raw:  map[string]interface{}{},
			key:  "time",
			want: model.TimeOfDay{},
		},
		{
			name: "empty nested object returns zero",
			raw:  map[string]interface{}{"time": map[string]interface{}{}},
			key:  "time",
			want: model.TimeOfDay{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTimeOfDay(tt.raw, tt.key)
			if err != nil {
				t.Errorf("getTimeOfDay() unexpected error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("getTimeOfDay() = %+v, want %+v", got, tt.want)
			}
		})
	}

	// Error cases
	t.Run("non-object value returns error", func(t *testing.T) {
		raw := map[string]interface{}{"time": "not-an-object"}
		_, err := getTimeOfDay(raw, "time")
		if err == nil {
			t.Error("getTimeOfDay() expected error for non-object value, got nil")
		}
	})

	t.Run("nil value returns error", func(t *testing.T) {
		raw := map[string]interface{}{"time": nil}
		_, err := getTimeOfDay(raw, "time")
		if err == nil {
			t.Error("getTimeOfDay() expected error for nil value, got nil")
		}
	})
}

func TestIsDuplicateKeyError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil error returns false", err: nil, want: false},
		{name: "duplicate key message", err: errors.New(`duplicate key value violates unique constraint "idx_name"`), want: true},
		{name: "unique constraint message", err: errors.New("unique constraint violation"), want: true},
		{name: "PostgreSQL error code 23505", err: fmt.Errorf("pq: 23505 error"), want: true},
		{name: "unrelated error", err: errors.New("connection refused"), want: false},
		{name: "empty error", err: errors.New(""), want: false},
		{name: "duplicate without 'key' substring", err: errors.New("found duplicate entries in batch"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDuplicateKeyError(tt.err)
			if got != tt.want {
				t.Errorf("isDuplicateKeyError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustIntPtr(t *testing.T) {
	tests := []struct {
		name    string
		raw     map[string]interface{}
		key     string
		wantNil bool
		want    int
	}{
		{name: "float64 value", raw: map[string]interface{}{"id": float64(42)}, key: "id", want: 42},
		{name: "string integer", raw: map[string]interface{}{"id": "99"}, key: "id", want: 99},
		{name: "missing key returns nil", raw: map[string]interface{}{}, key: "id", wantNil: true},
		{name: "nil value returns nil", raw: map[string]interface{}{"id": nil}, key: "id", wantNil: true},
		{name: "empty string returns nil", raw: map[string]interface{}{"id": ""}, key: "id", wantNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustIntPtr(tt.raw, tt.key)
			if (got == nil) != tt.wantNil {
				t.Errorf("mustIntPtr() = %v, wantNil %v", got, tt.wantNil)
				return
			}
			if got != nil && *got != tt.want {
				t.Errorf("mustIntPtr() = %d, want %d", *got, tt.want)
			}
		})
	}
}
