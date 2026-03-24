package loader

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type baseLoader struct {
	entityName string
	entityDir  string
	tableName  string
}

func newBaseLoader(entityName, entityDir, tableName string) baseLoader {
	return baseLoader{entityName: entityName, entityDir: entityDir, tableName: tableName}
}

func (l *baseLoader) EntityName() string {
	return l.entityName
}

func (l *baseLoader) loadFiles(
	ctx context.Context,
	db *bun.DB,
	datasetPath string,
	logger *zap.Logger,
	convertFn func(map[string]interface{}) (interface{}, error),
	batchSize int,
) (*LoadResult, error) {
	result := &LoadResult{Entity: l.entityName}

	pattern := filepath.Join(datasetPath, l.entityDir, "part-*.jsonl")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob %s: %w", pattern, err)
	}
	if len(files) == 0 {
		logger.Warn("no JSONL files found", zap.String("pattern", pattern))
		return result, nil
	}

	for _, file := range files {
		fileResult := l.loadFile(ctx, db, file, convertFn, batchSize)
		result.Merge(fileResult)
	}

	logger.Info("loaded entity",
		zap.String("entity", l.entityName),
		zap.Int("total", result.Total),
		zap.Int("inserted", result.Inserted),
		zap.Int("skipped", result.Skipped),
		zap.Int("errors", len(result.Errors)),
	)

	return result, nil
}

func (l *baseLoader) loadFile(
	ctx context.Context,
	db *bun.DB,
	file string,
	convertFn func(map[string]interface{}) (interface{}, error),
	batchSize int,
) *LoadResult {
	result := &LoadResult{Entity: l.entityName}

	f, err := os.Open(file)
	if err != nil {
		result.addError(0, file, "", "failed to open file", err.Error())
		return result
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	var batch []interface{}
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()

		var raw map[string]interface{}
		if err := json.Unmarshal(line, &raw); err != nil {
			result.addError(lineNum, file, "", "invalid JSON", string(line))
			result.Skipped++
			continue
		}

		record, convErr := convertFn(raw)
		if convErr != nil {
			result.addError(lineNum, file, convErr.(*ConvertError).Field, convErr.Error(), string(line))
			result.Skipped++
			continue
		}

		batch = append(batch, record)
		if len(batch) >= batchSize {
			inserted, skipped := l.insertBatch(ctx, db, batch, file, result)
			result.Inserted += inserted
			result.Skipped += skipped
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		inserted, skipped := l.insertBatch(ctx, db, batch, file, result)
		result.Inserted += inserted
		result.Skipped += skipped
	}

	result.Total = lineNum
	return result
}

// insertBatch uses bun's NewInsert with Returning("") to avoid the deadlock
// caused by bun.BaseModel auto-adding RETURNING "id".
func (l *baseLoader) insertBatch(ctx context.Context, db *bun.DB, batch []interface{}, file string, result *LoadResult) (inserted, skipped int) {
	if l.tableName == "" {
		result.addError(0, file, "", "unknown table name", "")
		return 0, len(batch)
	}

	typedSlice, err := toTypedSlice(batch)
	if err != nil {
		result.addError(0, file, "", "batch type conversion failed", err.Error())
		return 0, len(batch)
	}

	res, err := db.NewInsert().Model(typedSlice).On("CONFLICT DO NOTHING").Returning("").Exec(ctx)
	if err != nil {
		result.addError(0, file, "", "batch insert failed", err.Error())
		return 0, len(batch)
	}

	rowsAffected, _ := res.RowsAffected()
	return int(rowsAffected), len(batch) - int(rowsAffected)
}

func toTypedSlice(batch []interface{}) (interface{}, error) {
	if len(batch) == 0 {
		return nil, fmt.Errorf("empty batch")
	}
	sliceType := reflect.SliceOf(reflect.TypeOf(batch[0]))
	typedSlice := reflect.MakeSlice(sliceType, len(batch), len(batch))
	for i, record := range batch {
		typedSlice.Index(i).Set(reflect.ValueOf(record))
	}
	slicePtr := reflect.New(sliceType)
	slicePtr.Elem().Set(typedSlice)
	return slicePtr.Interface(), nil
}
