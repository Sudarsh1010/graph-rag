package api

import (
	"context"
	"encoding/json"
	"net/http"
)

type NLQueryServiceAdapter struct {
	queryFn func(ctx context.Context, req []byte) (int, interface{})
}

func (a *NLQueryServiceAdapter) QueryJSON(ctx context.Context, body []byte) (int, interface{}) {
	return a.queryFn(ctx, body)
}

func NewNLQueryAdapter(queryFn func(ctx context.Context, req []byte) (int, interface{})) NLQueryHandler {
	return &NLQueryServiceAdapter{queryFn: queryFn}
}

func NLQueryHandlerFunc(fn func(ctx context.Context, req NLQueryRequest) (*NLQueryResponse, error)) func(ctx context.Context, body []byte) (int, interface{}) {
	return func(ctx context.Context, body []byte) (int, interface{}) {
		var req NLQueryRequest
		if err := json.Unmarshal(body, &req); err != nil {
			return http.StatusBadRequest, map[string]string{"error": "invalid request body: " + err.Error()}
		}
		if req.Question == "" {
			return http.StatusBadRequest, map[string]string{"error": "question is required"}
		}

		resp, err := fn(ctx, req)
		if err != nil {
			return http.StatusInternalServerError, map[string]string{"error": err.Error()}
		}
		return http.StatusOK, resp
	}
}

type NLQueryRequest struct {
	Question       string `json:"question"`
	ConversationID string `json:"conversation_id,omitempty"`
}

type NLQueryResponse struct {
	Answer        string                   `json:"answer"`
	Data          []map[string]interface{} `json:"data,omitempty"`
	Visualization *VisualizationHint       `json:"visualization,omitempty"`
	Sources       []string                 `json:"sources"`
	QueryPlan     interface{}              `json:"query_plan,omitempty"`
}

type VisualizationHint struct {
	Type string `json:"type"`
	X    string `json:"x,omitempty"`
	Y    string `json:"y,omitempty"`
}
