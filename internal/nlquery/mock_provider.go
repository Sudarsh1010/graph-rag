package nlquery

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// MockProvider is a test LLMProvider that returns canned responses.
type MockProvider struct {
	mu       sync.Mutex
	queries  []MockQuery
	fallback string
}

// MockQuery maps a question substring to a canned response.
type MockQuery struct {
	Contains string // substring match on the user question
	Response string // JSON query plan to return
}

// NewMockProvider creates a MockProvider with a set of canned query mappings.
func NewMockProvider(queries []MockQuery, fallback string) *MockProvider {
	return &MockProvider{
		queries:  queries,
		fallback: fallback,
	}
}

// GenerateCompletion matches the user prompt against canned queries and returns
// the corresponding response. If no match is found, the fallback is used.
func (m *MockProvider) GenerateCompletion(_ context.Context, _, userPrompt string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, mq := range m.queries {
		if contains(userPrompt, mq.Contains) {
			return mq.Response, nil
		}
	}

	if m.fallback != "" {
		return m.fallback, nil
	}

	return "", fmt.Errorf("mock: no matching query for prompt: %s", truncate(userPrompt, 100))
}

// LastQueryPlan is a convenience method for tests to parse the fallback response
// as a QueryPlan.
func (m *MockProvider) LastQueryPlan() (*QueryPlan, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var plan QueryPlan
	if err := json.Unmarshal([]byte(m.fallback), &plan); err != nil {
		return nil, err
	}
	return &plan, nil
}

// DefaultMockProvider returns a MockProvider pre-loaded with common SAP query patterns.
func DefaultMockProvider() *MockProvider {
	return NewMockProvider([]MockQuery{
		{
			Contains: "order",
			Response: `{
  "steps": [
    {
      "type": "graph_query",
      "params": {
        "cypher": "MATCH (so:SalesOrder) RETURN so.salesOrder, so.totalNetAmount, so.transactionCurrency, so.overallDeliveryStatus LIMIT 10"
      }
    }
  ],
  "answer_format": "detailed_table"
}`,
		},
		{
			Contains: "customer",
			Response: `{
  "steps": [
    {
      "type": "graph_query",
      "params": {
        "cypher": "MATCH (c:Customer) RETURN c.businessPartner, c.businessPartnerFullName LIMIT 10"
      }
    }
  ],
  "answer_format": "list"
}`,
		},
		{
			Contains: "product",
			Response: `{
  "steps": [
    {
      "type": "graph_query",
      "params": {
        "cypher": "MATCH (p:Product) RETURN p.product, p.productGroup, p.division, p.baseUnit LIMIT 10"
      }
    }
  ],
  "answer_format": "detailed_table"
}`,
		},
	}, `{
  "steps": [
    {
      "type": "graph_query",
      "params": {
        "cypher": "MATCH (n) RETURN labels(n) AS type, count(n) AS count LIMIT 20"
      }
    }
  ],
  "answer_format": "brief_summary"
}`)
}

// contains is a case-insensitive substring match.
func contains(s, substr string) bool {
	return len(substr) > 0 && strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
