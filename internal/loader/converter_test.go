package loader

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, got time.Time)
	}{
		{
			name:  "RFC3339 format",
			input: "2024-03-15T10:30:00Z",
			check: func(t *testing.T, got time.Time) {
				if got.Year() != 2024 || got.Month() != time.March || got.Day() != 15 {
					t.Errorf("got %v, want date 2024-03-15", got)
				}
			},
		},
		{
			name:  "RFC3339 with timezone offset",
			input: "2024-03-15T10:30:00+05:30",
			check: func(t *testing.T, got time.Time) {
				if got.Year() != 2024 || got.Month() != time.March || got.Day() != 15 {
					t.Errorf("got %v, want date 2024-03-15", got)
				}
			},
		},
		{
			name:  "RFC3339Nano format",
			input: "2024-03-15T10:30:00.123456789Z",
			check: func(t *testing.T, got time.Time) {
				if got.Year() != 2024 {
					t.Errorf("got year %d, want 2024", got.Year())
				}
			},
		},
		{
			name:  "date only format",
			input: "2024-03-15",
			check: func(t *testing.T, got time.Time) {
				if got.Year() != 2024 || got.Month() != time.March || got.Day() != 15 {
					t.Errorf("got %v, want date 2024-03-15", got)
				}
			},
		},
		{
			name:  "datetime without timezone",
			input: "2024-03-15T10:30:00",
			check: func(t *testing.T, got time.Time) {
				if got.Hour() != 10 || got.Minute() != 30 {
					t.Errorf("got time %v, want 10:30", got)
				}
			},
		},
		{
			name:  "datetime with millis Z",
			input: "2024-03-15T10:30:00.000Z",
			check: func(t *testing.T, got time.Time) {
				if got.Year() != 2024 {
					t.Errorf("got year %d, want 2024", got.Year())
				}
			},
		},
		{
			name:  "datetime Z suffix",
			input: "2024-03-15T10:30:05Z",
			check: func(t *testing.T, got time.Time) {
				if got.Second() != 5 {
					t.Errorf("got second %d, want 5", got.Second())
				}
			},
		},
		{
			name:  "empty string returns zero time",
			input: "",
			check: func(t *testing.T, got time.Time) {
				if !got.IsZero() {
					t.Errorf("got %v, want zero time", got)
				}
			},
		},
		{
			name:  "whitespace string returns zero time",
			input: "   ",
			check: func(t *testing.T, got time.Time) {
				if !got.IsZero() {
					t.Errorf("got %v, want zero time", got)
				}
			},
		},
		{
			name:    "invalid format returns error",
			input:   "not-a-date",
			wantErr: true,
		},
		{
			name:    "random text returns error",
			input:   "hello world",
			wantErr: true,
		},
		{
			name:    "partially valid but wrong format",
			input:   "15/03/2024",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func TestParseTimePtr(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNil bool
		wantErr bool
	}{
		{name: "valid RFC3339", input: "2024-03-15T10:30:00Z", wantNil: false},
		{name: "valid date only", input: "2024-03-15", wantNil: false},
		{name: "empty string returns nil", input: "", wantNil: true},
		{name: "whitespace returns nil", input: "   ", wantNil: true},
		{name: "invalid format returns error", input: "bad-date", wantErr: true, wantNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTimePtr(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimePtr(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("ParseTimePtr(%q) = %v, wantNil %v", tt.input, got, tt.wantNil)
			}
			if got != nil && got.IsZero() {
				t.Errorf("ParseTimePtr(%q) returned zero time", tt.input)
			}
		})
	}
}

func TestParseStringPtr(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNil bool
		want    string
	}{
		{name: "non-empty string", input: "hello", want: "hello"},
		{name: "string with spaces", input: "hello world", want: "hello world"},
		{name: "empty string returns nil", input: "", wantNil: true},
		{name: "whitespace-only returns nil", input: "   ", wantNil: true},
		{name: "whitespace around text is trimmed", input: "  hello  ", want: "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseStringPtr(tt.input)
			if (got == nil) != tt.wantNil {
				t.Errorf("ParseStringPtr(%q) = %v, wantNil %v", tt.input, got, tt.wantNil)
				return
			}
			if got != nil && *got != tt.want {
				t.Errorf("ParseStringPtr(%q) = %q, want %q", tt.input, *got, tt.want)
			}
		})
	}
}

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		field   string
		wantErr bool
	}{
		{name: "non-empty value passes", value: "hello", field: "name"},
		{name: "empty value fails", value: "", field: "name", wantErr: true},
		{name: "whitespace-only value fails", value: "   ", field: "name", wantErr: true},
		{name: "value with leading/trailing spaces passes", value: "  hello  ", field: "name"},
		{name: "single character passes", value: "a", field: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequired(tt.value, tt.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequired(%q, %q) error = %v, wantErr %v", tt.value, tt.field, err, tt.wantErr)
			}
		})
	}
}

func TestParseDecimal(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    decimal.Decimal
		wantErr bool
	}{
		{name: "integer string", input: "42", want: decimal.RequireFromString("42")},
		{name: "float string", input: "3.14", want: decimal.RequireFromString("3.14")},
		{name: "negative number", input: "-10.5", want: decimal.RequireFromString("-10.5")},
		{name: "zero string", input: "0", want: decimal.RequireFromString("0")},
		{name: "scientific notation", input: "1.5e10", want: decimal.RequireFromString("1.5e10")},
		{name: "empty string returns zero", input: "", want: decimal.Decimal{}},
		{name: "whitespace returns zero", input: "   ", want: decimal.Decimal{}},
		{name: "spaces trimmed", input: "  3.14  ", want: decimal.RequireFromString("3.14")},
		{name: "invalid string returns error", input: "not-a-number", wantErr: true},
		{name: "letters returns error", input: "abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDecimal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDecimal(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ParseDecimal(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
