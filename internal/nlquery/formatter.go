package nlquery

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// formatResults sends the execution results to the LLM to produce a natural language
// answer to the user's original question.
func formatResults(
	ctx context.Context,
	llm LLMProvider,
	question string,
	results []map[string]interface{},
	answerFormat string,
) (string, error) {
	if len(results) == 0 {
		return "No results found for your query.", nil
	}

	// Serialize results for the LLM prompt. Limit to avoid token overflow.
	data := results
	if len(data) > 50 {
		data = data[:50]
	}

	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal results: %w", err)
	}

	systemPrompt := `You are a helpful assistant for SAP business data. Given the user's question and the query results below, provide a clear, concise answer. Use the answer_format hint to guide your response style. If the data is tabular, consider listing key rows. Always cite specific numbers and names from the data. Do not mention the query plan or technical details unless relevant.`

	userPrompt := fmt.Sprintf(
		"Question: %s\n\nAnswer format: %s\n\nQuery results (%d rows):\n```json\n%s\n```",
		question,
		answerFormat,
		len(data),
		string(dataJSON),
	)

	answer, err := llm.GenerateCompletion(ctx, systemPrompt, userPrompt)
	if err != nil {
		return "", fmt.Errorf("LLM formatting failed: %w", err)
	}

	return strings.TrimSpace(answer), nil
}
