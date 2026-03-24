package nlquery

import (
	"context"
	"fmt"
	"strings"

	"github.com/sudarsh1010/graph-rag/internal/api"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

// executePlan runs all steps in a query plan and returns merged results.
func executePlan(
	ctx context.Context,
	graph api.GraphQuerier,
	db *bun.DB,
	logger *zap.Logger,
	plan *QueryPlan,
) ([]map[string]interface{}, []string) {
	var allResults []map[string]interface{}
	var sources []string

	for i, step := range plan.Steps {
		logger.Info("executing step",
			zap.Int("index", i),
			zap.String("type", step.Type),
		)

		result, src, err := executeStep(ctx, graph, db, step)
		if err != nil {
			logger.Warn("step execution failed, skipping",
				zap.Int("index", i),
				zap.String("type", step.Type),
				zap.Error(err),
			)
			continue
		}

		allResults = append(allResults, result...)
		if src != "" {
			sources = append(sources, src)
		}
	}

	return allResults, dedup(sources)
}

// executeStep runs a single query step.
func executeStep(
	ctx context.Context,
	graph api.GraphQuerier,
	db *bun.DB,
	step QueryStep,
) ([]map[string]interface{}, string, error) {
	switch step.Type {
	case "graph_query":
		return executeGraphQuery(ctx, graph, step.Params)
	case "api_call":
		return executeSQLQuery(ctx, db, step.Params)
	case "aggregate":
		// aggregate is reserved for future server-side aggregation.
		// For now, fall through to graph_query if cypher is present.
		if cypher, ok := step.Params["cypher"]; ok {
			step.Params["cypher"] = cypher
			return executeGraphQuery(ctx, graph, step.Params)
		}
		return nil, "", fmt.Errorf("aggregate step requires 'cypher' param")
	default:
		return nil, "", fmt.Errorf("unknown step type: %s", step.Type)
	}
}

// executeGraphQuery runs a Cypher query against the AGE graph.
func executeGraphQuery(ctx context.Context, graph api.GraphQuerier, params map[string]string) ([]map[string]interface{}, string, error) {
	cypher := params["cypher"]
	if cypher == "" {
		return nil, "", fmt.Errorf("missing cypher param")
	}

	results, err := graph.Query(ctx, cypher)
	if err != nil {
		return nil, "", fmt.Errorf("graph query: %w", err)
	}

	return results, "graph:sap_ocg", nil
}

// executeSQLQuery runs a SQL query against the relational database.
func executeSQLQuery(ctx context.Context, db *bun.DB, params map[string]string) ([]map[string]interface{}, string, error) {
	sql := params["sql"]
	if sql == "" {
		return nil, "", fmt.Errorf("missing sql param")
	}

	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, "", fmt.Errorf("sql query: %w", err)
	}
	defer rows.Close()

	results, err := scanRows(rows)
	if err != nil {
		return nil, "", fmt.Errorf("scan sql results: %w", err)
	}

	// Derive source from the SQL query's table reference.
	source := extractTableFromSQL(sql)

	return results, source, nil
}

// scanRows reads all rows from sql.Rows into a slice of maps.
func scanRows(rows interface {
	Columns() ([]string, error)
	Next() bool
	Scan(...interface{}) error
	Err() error
}) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("get columns: %w", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		ptrs := make([]interface{}, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		row := make(map[string]interface{}, len(columns))
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}

// extractTableFromSQL attempts to extract the primary table name from a SQL query.
func extractTableFromSQL(sql string) string {
	upper := strings.ToUpper(strings.TrimSpace(sql))
	// Look for "FROM <table>" pattern.
	fromIdx := strings.Index(upper, " FROM ")
	if fromIdx == -1 {
		return "sql:unknown"
	}

	rest := strings.TrimSpace(sql[fromIdx+6:])
	// Extract the table name (stop at space, comma, or semicolon).
	for i, ch := range rest {
		if ch == ' ' || ch == ',' || ch == ';' || ch == '\n' {
			return "sql:" + rest[:i]
		}
	}
	return "sql:" + rest
}

// dedup removes duplicate strings while preserving order.
func dedup(ss []string) []string {
	seen := make(map[string]bool, len(ss))
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
