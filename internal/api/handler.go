package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/uptrace/bun"
)

// GraphQuerier is the interface for AGE graph queries.
type GraphQuerier interface {
	Query(ctx context.Context, cypher string) ([]map[string]interface{}, error)
}

type handler struct {
	db      *bun.DB
	graph   GraphQuerier
	nlQuery NLQueryHandler
}

type NLQueryHandler interface {
	QueryJSON(ctx context.Context, body []byte) (int, interface{})
}

func newHandler(db *bun.DB, graph GraphQuerier, nlq NLQueryHandler) *handler {
	return &handler{db: db, graph: graph, nlQuery: nlq}
}

func (h *handler) listEntities(w http.ResponseWriter, r *http.Request, table string) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")
	filter := r.URL.Query().Get("filter")

	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "20"
	}

	offset := (atoi(page) - 1) * atoi(limit)

	query := "SELECT * FROM " + sanitize(table)
	where := ""
	args := []interface{}{}

	if filter != "" {
		where = " WHERE " + filter
	}

	if sort != "" {
		where += " ORDER BY " + sanitize(sort)
	}

	query += fmt.Sprintf("%s LIMIT %s OFFSET %d", where, limit, offset)

	rows, err := h.db.QueryContext(r.Context(), query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query failed: "+err.Error())
		return
	}
	defer rows.Close()

	results, err := scanRows(rows)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "scan failed: "+err.Error())
		return
	}

	// Count total for pagination
	countQuery := "SELECT COUNT(*) FROM " + sanitize(table)
	var total int
	if err := h.db.QueryRowContext(r.Context(), countQuery).Scan(&total); err != nil {
		writeError(w, http.StatusInternalServerError, "count failed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  results,
		"page":  atoi(page),
		"limit": atoi(limit),
		"total": total,
	})
}

func (h *handler) getEntity(w http.ResponseWriter, r *http.Request, table, idColumn string) {
	id := chiURLParam(r, "id")

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 LIMIT 1", sanitize(table), sanitize(idColumn))

	rows, err := h.db.QueryContext(r.Context(), query, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query failed: "+err.Error())
		return
	}
	defer rows.Close()

	results, err := scanRows(rows)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "scan failed: "+err.Error())
		return
	}

	if len(results) == 0 {
		writeError(w, http.StatusNotFound, "entity not found")
		return
	}

	writeJSON(w, http.StatusOK, results[0])
}

func (h *handler) graphQuery(w http.ResponseWriter, r *http.Request) {
	if h.graph == nil {
		writeError(w, http.StatusServiceUnavailable, "graph service not available")
		return
	}

	var req struct {
		Cypher string `json:"cypher"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Cypher) == "" {
		writeError(w, http.StatusBadRequest, "cypher query is required")
		return
	}

	// Basic validation — only allow SELECT/MATCH (no CREATE/DELETE/DROP)
	upper := strings.ToUpper(strings.TrimSpace(req.Cypher))
	if strings.Contains(upper, "CREATE ") || strings.Contains(upper, "DELETE ") ||
		strings.Contains(upper, "DROP ") || strings.Contains(upper, "SET ") {
		writeError(w, http.StatusForbidden, "only read-only Cypher queries are allowed")
		return
	}

	results, err := h.graph.Query(r.Context(), req.Cypher)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "graph query failed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": results,
	})
}

func (h *handler) nlQueryHandler(w http.ResponseWriter, r *http.Request) {
	if h.nlQuery == nil {
		writeError(w, http.StatusServiceUnavailable, "NL query service not available")
		return
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	status, data := h.nlQuery.QueryJSON(r.Context(), buf.Bytes())
	writeJSON(w, status, data)
}

func (h *handler) searchEntities(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	entityType := r.URL.Query().Get("type")

	if query == "" {
		writeError(w, http.StatusBadRequest, "query parameter is required")
		return
	}

	// For now, perform basic ILIKE search across common text fields.
	// Vector search will be added via the embeddings pipeline.
	like := "%" + query + "%"
	searchableTables := map[string][]string{
		"products":          {"product", "product_group", "division"},
		"customers":         {"business_partner", "full_name"},
		"plants":            {"plant", "name"},
		"sales-orders":      {"sales_order", "sales_order_type"},
		"billing-documents": {"billing_document", "billing_document_type"},
		"deliveries":        {"delivery_document", "delivery_document_type"},
		"":                  {}, // search all
	}

	results := map[string][]map[string]interface{}{}
	tablesToSearch := searchableTables[entityType]
	if entityType == "" {
		// Search all tables
		for table, _ := range searchableTables {
			tablesToSearch = append(tablesToSearch, table)
		}
	}

	// Map frontend entity names to DB table names
	tableMap := map[string]string{
		"products":          "products",
		"customers":         "business_partners",
		"plants":            "plants",
		"sales-orders":      "sales_order_headers",
		"billing-documents": "billing_document_headers",
		"deliveries":        "outbound_delivery_headers",
	}

	for _, et := range []string{"products", "customers", "plants", "sales-orders", "billing-documents", "deliveries"} {
		if entityType != "" && entityType != et {
			continue
		}

		tableName, ok := tableMap[et]
		if !ok {
			continue
		}

		cols := searchableTables[et]
		if len(cols) == 0 {
			continue
		}

		conditions := []string{}
		for _, col := range cols {
			conditions = append(conditions, fmt.Sprintf("CAST(%s AS TEXT) ILIKE $1", sanitize(col)))
		}

		whereClause := strings.Join(conditions, " OR ")
		q := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 10", sanitize(tableName), whereClause)

		rows, err := h.db.QueryContext(r.Context(), q, like)
		if err != nil {
			continue
		}

		scanned, err := scanRows(rows)
		rows.Close()
		if err != nil {
			continue
		}

		if len(scanned) > 0 {
			results[et] = scanned
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"query":   query,
		"results": results,
	})
}

// scanRows reads all rows from a sql.Rows and returns a slice of maps.
func scanRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	defer rows.Close()

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
			// Convert []byte to string for JSON serialization
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

// sanitize prevents SQL injection by rejecting identifiers with special chars.
func sanitize(name string) string {
	// Only allow alphanumeric and underscores
	safe := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return -1
	}, name)
	return safe
}

// chiURLParam wraps chi.URLParam to avoid importing chi in handler tests.
var chiURLParam = func(r *http.Request, key string) string {
	return chiURLParamImpl(r, key)
}

func chiURLParamImpl(r *http.Request, key string) string {
	return r.PathValue(key)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]interface{}{
		"error": map[string]interface{}{
			"code":    status,
			"message": message,
		},
	})
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
