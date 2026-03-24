package nlquery

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// translate sends the user question to the LLM with the system prompt and parses
// the structured JSON QueryPlan from the response.
func translate(ctx context.Context, llm LLMProvider, question string) (*QueryPlan, error) {
	systemPrompt := BuildSystemPrompt()

	userPrompt := fmt.Sprintf("Convert this question into a query plan:\n\n%s", question)

	raw, err := llm.GenerateCompletion(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("LLM completion failed: %w", err)
	}

	plan, err := parseQueryPlan(raw)
	if err != nil {
		return nil, fmt.Errorf("parse query plan: %w", err)
	}

	if err := validatePlan(plan); err != nil {
		return nil, fmt.Errorf("validate query plan: %w", err)
	}

	return plan, nil
}

// parseQueryPlan extracts and deserializes the JSON query plan from the LLM response.
// It handles cases where the LLM wraps the JSON in markdown code fences.
func parseQueryPlan(raw string) (*QueryPlan, error) {
	trimmed := strings.TrimSpace(raw)

	// Strip markdown code fences if present.
	if strings.HasPrefix(trimmed, "```json") {
		trimmed = strings.TrimPrefix(trimmed, "```json")
		trimmed = strings.TrimSuffix(trimmed, "```")
		trimmed = strings.TrimSpace(trimmed)
	} else if strings.HasPrefix(trimmed, "```") {
		trimmed = strings.TrimPrefix(trimmed, "```")
		trimmed = strings.TrimSuffix(trimmed, "```")
		trimmed = strings.TrimSpace(trimmed)
	}

	var plan QueryPlan
	if err := json.Unmarshal([]byte(trimmed), &plan); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w (raw: %s)", err, truncate(raw, 200))
	}

	if len(plan.Steps) == 0 {
		return nil, fmt.Errorf("query plan must contain at least one step")
	}

	return &plan, nil
}

// validatePlan checks that all query steps are read-only and have required params.
func validatePlan(plan *QueryPlan) error {
	for i, step := range plan.Steps {
		switch step.Type {
		case "graph_query":
			cypher, ok := step.Params["cypher"]
			if !ok || strings.TrimSpace(cypher) == "" {
				return fmt.Errorf("step %d (graph_query): missing 'cypher' param", i)
			}
			if err := validateReadOnlyCypher(cypher); err != nil {
				return fmt.Errorf("step %d (graph_query): %w", i, err)
			}

		case "api_call":
			sql, ok := step.Params["sql"]
			if !ok || strings.TrimSpace(sql) == "" {
				return fmt.Errorf("step %d (api_call): missing 'sql' param", i)
			}
			if err := validateReadOnlySQL(sql); err != nil {
				return fmt.Errorf("step %d (api_call): %w", i, err)
			}

		case "aggregate":
			// aggregate steps are reserved for future use; pass through.

		default:
			return fmt.Errorf("step %d: unknown step type %q", i, step.Type)
		}
	}

	return nil
}

// validateReadOnlyCypher rejects Cypher queries containing write operations.
func validateReadOnlyCypher(cypher string) error {
	upper := strings.ToUpper(strings.TrimSpace(cypher))
	forbidden := []string{"CREATE ", "DELETE ", "DETACH DELETE", "DROP ", "SET ", "MERGE ", "REMOVE "}
	for _, op := range forbidden {
		if strings.Contains(upper, op) {
			return fmt.Errorf("write operation %q is not allowed", strings.TrimSpace(op))
		}
	}
	return nil
}

// validateReadOnlySQL rejects SQL queries containing write operations.
func validateReadOnlySQL(sql string) error {
	upper := strings.ToUpper(strings.TrimSpace(sql))
	// Only allow SELECT statements.
	if !strings.HasPrefix(upper, "SELECT") {
		return fmt.Errorf("only SELECT statements are allowed")
	}
	forbidden := []string{"INSERT ", "UPDATE ", "DELETE ", "DROP ", "ALTER ", "CREATE ", "GRANT ", "TRUNCATE "}
	for _, op := range forbidden {
		if strings.Contains(upper, op) {
			return fmt.Errorf("write operation %q is not allowed", strings.TrimSpace(op))
		}
	}
	return nil
}

// planToVisualization extracts a VisualizationHint from the query plan if present
// in the LLM response, or returns nil.
func planToVisualization(plan *QueryPlan) *VisualizationHint {
	if plan == nil || len(plan.Steps) == 0 {
		return nil
	}

	// Check if the first step has visualization params embedded.
	// The LLM may include these in the plan. We look for common patterns.
	// By default, return a simple hint.
	if plan.AnswerFormat == "detailed_table" {
		return &VisualizationHint{Type: "table"}
	}

	return &VisualizationHint{Type: "none"}
}

// truncate shortens a string to maxLen for error messages.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
