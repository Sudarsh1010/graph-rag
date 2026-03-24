package loader

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type LoadError struct {
	Line    int
	File    string
	Field   string
	Message string
	Raw     string
}

type LoadResult struct {
	Entity   string
	Total    int
	Inserted int
	Skipped  int
	Errors   []LoadError
}

func (r *LoadResult) addError(line int, file, field, message, raw string) {
	trimmed := raw
	if len(trimmed) > 200 {
		trimmed = trimmed[:200]
	}
	r.Errors = append(r.Errors, LoadError{
		Line:    line,
		File:    file,
		Field:   field,
		Message: message,
		Raw:     trimmed,
	})
}

func (r *LoadResult) Merge(other *LoadResult) {
	r.Total += other.Total
	r.Inserted += other.Inserted
	r.Skipped += other.Skipped
	r.Errors = append(r.Errors, other.Errors...)
}

type EntityLoader interface {
	EntityName() string
	Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error)
}

type ConvertError struct {
	Field   string
	Message string
}

func (e *ConvertError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func newConvertError(field, message string) *ConvertError {
	return &ConvertError{Field: field, Message: message}
}

func ParseDecimal(s string) (decimal.Decimal, error) {
	s = trimSpaces(s)
	if s == "" {
		return decimal.Decimal{}, nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("invalid decimal %q: %w", s, err)
	}
	return d, nil
}

func ParseRequiredDecimal(s string, field string) (decimal.Decimal, error) {
	s = trimSpaces(s)
	if s == "" {
		return decimal.Decimal{}, newConvertError(field, "required decimal is empty")
	}
	return ParseDecimal(s)
}

func ParseNonNegDecimal(s string, field string) (decimal.Decimal, error) {
	d, err := ParseDecimal(s)
	if err != nil {
		return d, err
	}
	if d.IsNegative() {
		return decimal.Decimal{}, newConvertError(field, fmt.Sprintf("value must be non-negative, got %s", d.String()))
	}
	return d, nil
}
