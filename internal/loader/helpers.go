package loader

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sudarsh1010/graph-rag/internal/model"
)

func trimSpaces(s string) string {
	return strings.TrimSpace(s)
}

func getString(raw map[string]interface{}, key string) string {
	v, ok := raw[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	return s
}

func getStringPtr(raw map[string]interface{}, key string) *string {
	v, ok := raw[key]
	if !ok {
		return nil
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return nil
	}
	return &s
}

func getDecimal(raw map[string]interface{}, key string) (decimal.Decimal, error) {
	s := getString(raw, key)
	return ParseDecimal(s)
}

func getTime(raw map[string]interface{}, key string) (time.Time, error) {
	s := getString(raw, key)
	return ParseTime(s)
}

func getTimePtr(raw map[string]interface{}, key string) (*time.Time, error) {
	v, ok := raw[key]
	if !ok {
		return nil, nil
	}

	switch val := v.(type) {
	case nil:
		return nil, nil
	case string:
		if val == "" {
			return nil, nil
		}
		t, err := ParseTime(val)
		if err != nil {
			return nil, newConvertError(key, err.Error())
		}
		return &t, nil
	default:
		return nil, newConvertError(key, fmt.Sprintf("expected string or null, got %T", v))
	}
}

func getBool(raw map[string]interface{}, key string) bool {
	v, ok := raw[key]
	if !ok {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return strings.ToLower(val) == "true"
	default:
		return false
	}
}

func getIntRequired(raw map[string]interface{}, key string) (int, error) {
	v, ok := raw[key]
	if !ok {
		return 0, newConvertError(key, "missing required field")
	}
	switch val := v.(type) {
	case float64:
		return int(val), nil
	case string:
		var i int
		_, err := fmt.Sscanf(val, "%d", &i)
		if err != nil {
			return 0, newConvertError(key, fmt.Sprintf("invalid integer: %q", val))
		}
		return i, nil
	case int:
		return val, nil
	default:
		return 0, newConvertError(key, fmt.Sprintf("expected number, got %T", v))
	}
}

func getIntPtr(raw map[string]interface{}, key string) (*int, error) {
	v, ok := raw[key]
	if !ok {
		return nil, nil
	}
	switch val := v.(type) {
	case nil:
		return nil, nil
	case float64:
		i := int(val)
		return &i, nil
	case string:
		if val == "" {
			return nil, nil
		}
		var i int
		_, err := fmt.Sscanf(val, "%d", &i)
		if err != nil {
			return nil, newConvertError(key, fmt.Sprintf("invalid integer: %q", val))
		}
		return &i, nil
	default:
		return nil, newConvertError(key, fmt.Sprintf("expected number or null, got %T", v))
	}
}

func getTimeOfDay(raw map[string]interface{}, key string) (model.TimeOfDay, error) {
	v, ok := raw[key]
	if !ok {
		return model.TimeOfDay{}, nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return model.TimeOfDay{}, newConvertError(key, fmt.Sprintf("expected object, got %T", v))
	}

	var tod model.TimeOfDay
	if h, err := getIntFromMap(m, "hours"); err == nil {
		tod.Hours = h
	}
	if min, err := getIntFromMap(m, "minutes"); err == nil {
		tod.Minutes = min
	}
	if sec, err := getIntFromMap(m, "seconds"); err == nil {
		tod.Seconds = sec
	}
	return tod, nil
}

func getIntFromMap(m map[string]interface{}, key string) (int, error) {
	v, ok := m[key]
	if !ok {
		return 0, nil
	}
	switch val := v.(type) {
	case float64:
		return int(val), nil
	case int:
		return val, nil
	default:
		return 0, fmt.Errorf("expected number for %s, got %T", key, v)
	}
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "23505")
}

func mustIntPtr(raw map[string]interface{}, key string) *int {
	v, _ := getIntPtr(raw, key)
	return v
}
