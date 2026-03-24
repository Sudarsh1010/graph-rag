package store

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type VectorSearchResult struct {
	ID         string
	Similarity float64
}

func SearchByVector(ctx context.Context, db *bun.DB, table, vectorColumn, queryVector string, limit int) ([]VectorSearchResult, error) {
	query := fmt.Sprintf(
		`SELECT %s as id, 1 - (embedding <=> %s::vector) as similarity
		 FROM %s
		 WHERE embedding IS NOT NULL
		 ORDER BY embedding <=> %s::vector
		 LIMIT %d`,
		table, queryVector, table, queryVector, limit,
	)

	var results []VectorSearchResult
	err := db.NewRaw(query).Scan(ctx, &results)
	if err != nil {
		return nil, fmt.Errorf("vector search on %s: %w", table, err)
	}

	return results, nil
}
