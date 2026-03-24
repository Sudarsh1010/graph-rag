package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/uptrace/bun"
)

type EmbeddingRecord struct {
	ID        string    `json:"id"`
	Embedding []float32 `json:"embedding"`
}

func LoadEmbeddings(ctx context.Context, db *bun.DB, embeddingsPath string) error {
	files, err := filepath.Glob(filepath.Join(embeddingsPath, "*.jsonl"))
	if err != nil {
		return fmt.Errorf("glob embeddings: %w", err)
	}

	for _, file := range files {
		baseName := strings.TrimSuffix(filepath.Base(file), "_embeddings.jsonl")
		entityType := strings.Replace(baseName, "_", "", -1)

		count, err := loadEmbeddingFile(ctx, db, file, entityType)
		if err != nil {
			return fmt.Errorf("load %s: %w", file, err)
		}
		fmt.Printf("Loaded %d embeddings for %s\n", count, entityType)
	}

	return nil
}

func loadEmbeddingFile(ctx context.Context, db *bun.DB, file string, entityType string) (int, error) {
	f, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	count := 0

	for decoder.More() {
		var rec EmbeddingRecord
		if err := decoder.Decode(&rec); err != nil {
			return count, fmt.Errorf("decode: %w", err)
		}

		if len(rec.Embedding) == 0 {
			continue
		}

		vectorStr := formatVector(rec.Embedding)

		var table, idColumn string
		switch entityType {
		case "productdescriptions":
			table = "product_descriptions"
			idColumn = "product"
		case "businesspartners":
			table = "business_partners"
			idColumn = "business_partner"
		case "plants":
			table = "plants"
			idColumn = "plant"
		default:
			continue
		}

		parts := strings.SplitN(rec.ID, "_", 2)
		id := parts[0]

		var subID string
		if len(parts) > 1 {
			subID = parts[1]
		}

		var query string
		var args []interface{}

		if subID != "" {
			query = fmt.Sprintf(
				"UPDATE %s SET embedding = %s::vector WHERE %s = ? AND language = ?",
				table, vectorStr, idColumn,
			)
			args = []interface{}{id, subID}
		} else {
			query = fmt.Sprintf(
				"UPDATE %s SET embedding = %s::vector WHERE %s = ?",
				table, vectorStr, idColumn,
			)
			args = []interface{}{id}
		}

		_, err := db.ExecContext(ctx, query, args...)
		if err != nil {
			return count, fmt.Errorf("update %s id=%s: %w", table, id, err)
		}
		count++
	}

	return count, nil
}

func formatVector(embedding []float32) string {
	strs := make([]string, len(embedding))
	for i, v := range embedding {
		strs[i] = fmt.Sprintf("%g", float64(v))
	}
	return "[" + strings.Join(strs, ",") + "]"
}
