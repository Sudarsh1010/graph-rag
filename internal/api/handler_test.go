package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSanitize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple alphanumeric", input: "users", want: "users"},
		{name: "with underscores", input: "user_accounts", want: "user_accounts"},
		{name: "with numbers", input: "table123", want: "table123"},
		{name: "uppercase letters", input: "TableName", want: "TableName"},
		{name: "SQL injection single quote", input: "users'; DROP TABLE users;--", want: "usersDROPTABLEusers"},
		{name: "SQL injection double dash", input: "users--drop", want: "usersdrop"},
		{name: "spaces are stripped", input: "my table", want: "mytable"},
		{name: "special characters stripped", input: "col@#$%name", want: "colname"},
		{name: "semicolon stripped", input: "table;drop", want: "tabledrop"},
		{name: "parens stripped", input: "col(1)", want: "col1"},
		{name: "dots stripped", input: "schema.table", want: "schematable"},
		{name: "empty string", input: "", want: ""},
		{name: "all special chars", input: "!@#$%^&*()", want: ""},
		{name: "mixed valid and invalid", input: "a1_b2-c3@d4", want: "a1_b2c3d4"},
		{name: "newlines stripped", input: "col\nname", want: "colname"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitize(tt.input)
			if got != tt.want {
				t.Errorf("sanitize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAtoi(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{name: "valid positive integer", input: "42", want: 42},
		{name: "zero", input: "0", want: 0},
		{name: "negative integer", input: "-5", want: -5},
		{name: "large number", input: "999999", want: 999999},
		{name: "empty string defaults to 0", input: "", want: 0},
		{name: "non-numeric defaults to 0", input: "abc", want: 0},
		{name: "float string truncated to 0", input: "3.14", want: 0},
		{name: "mixed alphanumeric defaults to 0", input: "12abc", want: 0},
		{name: "whitespace defaults to 0", input: "  ", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := atoi(tt.input)
			if got != tt.want {
				t.Errorf("atoi(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		data     interface{}
		wantCode int
		check    func(t *testing.T, body []byte)
	}{
		{
			name:     "200 OK with map data",
			status:   http.StatusOK,
			data:     map[string]interface{}{"key": "value"},
			wantCode: http.StatusOK,
			check: func(t *testing.T, body []byte) {
				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("failed to unmarshal JSON: %v", err)
				}
				if result["key"] != "value" {
					t.Errorf("got key=%v, want value", result["key"])
				}
			},
		},
		{
			name:     "201 Created with string slice",
			status:   http.StatusCreated,
			data:     []string{"a", "b", "c"},
			wantCode: http.StatusCreated,
			check: func(t *testing.T, body []byte) {
				var result []string
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("failed to unmarshal JSON: %v", err)
				}
				if len(result) != 3 {
					t.Errorf("got %d items, want 3", len(result))
				}
			},
		},
		{
			name:     "500 with nil data",
			status:   http.StatusInternalServerError,
			data:     nil,
			wantCode: http.StatusInternalServerError,
			check: func(t *testing.T, body []byte) {
				if strings.TrimSpace(string(body)) != "null" {
					t.Errorf("got body %q, want null", string(body))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeJSON(rec, tt.status, tt.data)

			if rec.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantCode)
			}

			ct := rec.Header().Get("Content-Type")
			if !strings.Contains(ct, "application/json") {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}

			if tt.check != nil {
				tt.check(t, rec.Body.Bytes())
			}
		})
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		message  string
		wantCode int
		check    func(t *testing.T, body []byte)
	}{
		{
			name:     "400 bad request",
			status:   http.StatusBadRequest,
			message:  "invalid input",
			wantCode: http.StatusBadRequest,
			check: func(t *testing.T, body []byte) {
				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("failed to unmarshal JSON: %v", err)
				}
				errObj, ok := result["error"].(map[string]interface{})
				if !ok {
					t.Fatal("response.error is not an object")
				}
				if errObj["code"] != float64(http.StatusBadRequest) {
					t.Errorf("error.code = %v, want %d", errObj["code"], http.StatusBadRequest)
				}
				if errObj["message"] != "invalid input" {
					t.Errorf("error.message = %v, want 'invalid input'", errObj["message"])
				}
			},
		},
		{
			name:     "404 not found",
			status:   http.StatusNotFound,
			message:  "entity not found",
			wantCode: http.StatusNotFound,
			check: func(t *testing.T, body []byte) {
				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("failed to unmarshal JSON: %v", err)
				}
				errObj, ok := result["error"].(map[string]interface{})
				if !ok {
					t.Fatal("response.error is not an object")
				}
				if errObj["code"] != float64(http.StatusNotFound) {
					t.Errorf("error.code = %v, want %d", errObj["code"], http.StatusNotFound)
				}
			},
		},
		{
			name:     "500 internal server error",
			status:   http.StatusInternalServerError,
			message:  "something went wrong",
			wantCode: http.StatusInternalServerError,
			check: func(t *testing.T, body []byte) {
				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("failed to unmarshal JSON: %v", err)
				}
				errObj, ok := result["error"].(map[string]interface{})
				if !ok {
					t.Fatal("response.error is not an object")
				}
				if errObj["code"] != float64(http.StatusInternalServerError) {
					t.Errorf("error.code = %v, want %d", errObj["code"], http.StatusInternalServerError)
				}
				if errObj["message"] != "something went wrong" {
					t.Errorf("error.message = %v, want 'something went wrong'", errObj["message"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeError(rec, tt.status, tt.message)

			if rec.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantCode)
			}

			ct := rec.Header().Get("Content-Type")
			if !strings.Contains(ct, "application/json") {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}

			if tt.check != nil {
				tt.check(t, rec.Body.Bytes())
			}
		})
	}
}
