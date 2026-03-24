package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

const graphName = "sap_ocg"
const ageSearchPath = "SET LOCAL search_path TO ag_catalog, public"

type Service struct {
	db     *bun.DB
	logger *zap.Logger
}

func NewService(db *bun.DB, logger *zap.Logger) *Service {
	return &Service{db: db, logger: logger}
}

func (s *Service) Initialize(ctx context.Context) error {
	s.logger.Info("initializing AGE graph", zap.String("graph", graphName))

	var exists bool
	err := s.db.NewRaw("SELECT EXISTS(SELECT 1 FROM ag_catalog.ag_graph WHERE name = ?)", graphName).
		Scan(ctx, &exists)
	if err != nil {
		return fmt.Errorf("check graph exists: %w", err)
	}

	if exists {
		s.logger.Info("graph already exists, skipping creation", zap.String("graph", graphName))
		return nil
	}

	_, err = s.db.ExecContext(ctx, "SELECT ag_catalog.create_graph(?)", graphName)
	if err != nil {
		return fmt.Errorf("create graph: %w", err)
	}

	s.logger.Info("graph created successfully", zap.String("graph", graphName))
	return nil
}

func (s *Service) LoadVertices(ctx context.Context) error {
	for _, vd := range vertexDefs {
		count, err := s.loadVertexType(ctx, vd)
		if err != nil {
			return fmt.Errorf("load %s vertices: %w", vd.label, err)
		}
		s.logger.Info("loaded vertices", zap.String("label", vd.label), zap.Int("count", count))
	}
	return nil
}

func (s *Service) loadVertexType(ctx context.Context, vd vertexDef) (int, error) {
	cols := strings.Join(vd.columns, ", ")
	query := fmt.Sprintf("SELECT %s FROM %s", cols, vd.table)

	var rows []map[string]interface{}
	scanErr := s.db.NewRaw(query).Scan(ctx, &rows)
	if scanErr != nil {
		return 0, fmt.Errorf("query %s: %w", vd.table, scanErr)
	}
	if len(rows) == 0 {
		return 0, nil
	}

	txErr := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.ExecContext(ctx, ageSearchPath); err != nil {
			return fmt.Errorf("set search_path: %w", err)
		}
		batchSize := 500
		for i := 0; i < len(rows); i += batchSize {
			end := i + batchSize
			if end > len(rows) {
				end = len(rows)
			}
			batch := rows[i:end]

			agtypeList, listErr := buildAgtypeList(batch, vd.columns)
			if listErr != nil {
				return listErr
			}

			cypher := fmt.Sprintf("WITH %s AS batch UNWIND batch AS row CREATE (n:%s {%s})",
				agtypeList, vd.label, buildPropertyMap(vd.columns, "row"))

			q := fmt.Sprintf("SELECT * FROM cypher('%s', $$ %s $$) AS (n agtype)", graphName, cypher)
			if _, err := tx.ExecContext(ctx, q); err != nil {
				return fmt.Errorf("create %s batch [%d:%d]: %w", vd.label, i, end, err)
			}
		}
		return nil
	})

	if txErr != nil {
		return 0, txErr
	}
	return len(rows), nil
}

func buildAgtypeList(rows []map[string]interface{}, columns []string) (string, error) {
	items := make([]string, len(rows))
	for i, r := range rows {
		props := make([]string, len(columns))
		for j, col := range columns {
			props[j] = fmt.Sprintf("%s: %s", col, agtypeValue(r[col]))
		}
		items[i] = "{" + strings.Join(props, ", ") + "}"
	}
	return "[" + strings.Join(items, ", ") + "]", nil
}

// agtypeValue formats a Go value as an AGE Cypher literal.
// Strings are single-quoted (with escaping), numbers are passed through,
// nil becomes null.
func agtypeValue(v interface{}) string {
	if v == nil {
		return "null"
	}
	switch val := v.(type) {
	case string:
		// Escape single quotes for Cypher string literals
		escaped := strings.ReplaceAll(val, "'", "\\'")
		return fmt.Sprintf("'%s'", escaped)
	case float64:
		// Preserve as number (no quotes)
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		// Fallback: try JSON then wrap as string
		b, err := json.Marshal(val)
		if err != nil {
			return "null"
		}
		return string(b)
	}
}

func buildPropertyMap(columns []string, prefix string) string {
	props := make([]string, len(columns))
	for i, col := range columns {
		props[i] = fmt.Sprintf("%s: %s.%s", col, prefix, col)
	}
	return strings.Join(props, ", ")
}

func (s *Service) LoadEdges(ctx context.Context) error {
	for _, ed := range edgeDefs {
		count, err := s.loadEdgeType(ctx, ed)
		if err != nil {
			return fmt.Errorf("load %s edges: %w", ed.label, err)
		}
		s.logger.Info("loaded edges", zap.String("label", ed.label), zap.Int("count", count))
	}
	return nil
}

func (s *Service) loadEdgeType(ctx context.Context, ed edgeDef) (int, error) {
	createCypher := fmt.Sprintf(
		"MATCH (a:%s), (b:%s) WHERE %s CREATE (a)-[:%s]->(b)",
		ed.from, ed.to, ed.match, ed.label,
	)

	var count int
	txErr := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.ExecContext(ctx, ageSearchPath); err != nil {
			return fmt.Errorf("set search_path: %w", err)
		}
		q := fmt.Sprintf("SELECT * FROM cypher('%s', $$ %s $$) AS (e agtype)", graphName, createCypher)
		if _, err := tx.ExecContext(ctx, q); err != nil {
			return fmt.Errorf("create edges: %w", err)
		}
		countCypher := fmt.Sprintf("MATCH ()-[r:%s]->() RETURN count(r)", ed.label)
		cq := fmt.Sprintf("SELECT * FROM cypher('%s', $$ %s $$) AS (cnt bigint)", graphName, countCypher)
		return tx.NewRaw(cq).Scan(ctx, &count)
	})

	if txErr != nil {
		return 0, txErr
	}
	return count, nil
}

func (s *Service) FullLoad(ctx context.Context) error {
	if err := s.Initialize(ctx); err != nil {
		return err
	}
	if err := s.LoadVertices(ctx); err != nil {
		return err
	}
	return s.LoadEdges(ctx)
}

func (s *Service) Query(ctx context.Context, cypher string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.ExecContext(ctx, ageSearchPath); err != nil {
			return fmt.Errorf("set search_path: %w", err)
		}
		q := fmt.Sprintf("SELECT * FROM cypher('%s', $$ %s $$) AS (result agtype)", graphName, cypher)
		return tx.NewRaw(q).Scan(ctx, &results)
	})
	if err != nil {
		return nil, fmt.Errorf("graph query: %w", err)
	}
	return results, nil
}
