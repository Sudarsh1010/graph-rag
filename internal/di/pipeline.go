package di

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type PipelineStage string

const (
	StageMigrations     PipelineStage = "migrations"
	StageJSONLLoad      PipelineStage = "jsonl_load"
	StageEmbeddingsLoad PipelineStage = "embeddings_load"
	StageGraphVertices  PipelineStage = "graph_vertices"
	StageGraphEdges     PipelineStage = "graph_edges"
)

type PipelineService struct {
	db     *bun.DB
	logger *zap.Logger
}

func NewPipelineService(db *Database, logger *zap.Logger) *PipelineService {
	return &PipelineService{db: db.DB, logger: logger}
}

func (ps *PipelineService) Init(ctx context.Context) error {
	// Safety net: ensure pipeline_stages table exists (independent of migration)
	query := `
		CREATE TABLE IF NOT EXISTS pipeline_stages (
			stage VARCHAR(50) PRIMARY KEY,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			started_at TIMESTAMPTZ,
			completed_at TIMESTAMPTZ,
			error_message TEXT,
			metadata JSONB DEFAULT '{}'
		);`
	_, err := ps.db.ExecContext(ctx, query)
	return err
}

func (ps *PipelineService) RunStage(ctx context.Context, stage PipelineStage, fn func(context.Context) error) error {
	startTime := time.Now()

	// Get current status
	var status string
	var errorMessage string
	err := ps.db.NewRaw(`SELECT status, error_message FROM pipeline_stages WHERE stage = ?`, stage).Scan(ctx, &status, &errorMessage)
	if err != nil {
		// If no row exists, insert with pending status
		if err.Error() == "sql: no rows in result set" {
			status = "pending"
			errorMessage = ""
			_, err = ps.db.ExecContext(ctx, `INSERT INTO pipeline_stages (stage, status) VALUES (?, 'pending')`, stage)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Handle existing status
	switch status {
	case "completed":
		ps.logger.Info("stage " + string(stage) + ": skipping (already completed)")
		return nil
	case "failed", "running":
		ps.logger.Info("stage " + string(stage) + ": retrying (previous status: " + status + ")")
	default:
		// pending or any other status - proceed
	}

	// Set status to running
	_, err = ps.db.ExecContext(ctx, `UPDATE pipeline_stages SET status = 'running', started_at = ?, error_message = NULL WHERE stage = ?`, time.Now(), stage)
	if err != nil {
		return err
	}

	// Execute the function
	if err = fn(ctx); err != nil {
		// On failure
		ps.logger.Error("stage "+string(stage)+": failed", zap.Error(err))
		_, updateErr := ps.db.ExecContext(ctx, `UPDATE pipeline_stages SET status = 'failed', error_message = ? WHERE stage = ?`, err.Error(), stage)
		if updateErr != nil {
			return updateErr
		}
		return err
	}

	// On success
	ps.logger.Info("stage "+string(stage)+": completed", zap.Duration("elapsed", time.Since(startTime)))
	_, err = ps.db.ExecContext(ctx, `UPDATE pipeline_stages SET status = 'completed', completed_at = ? WHERE stage = ?`, time.Now(), stage)
	return err
}

func (ps *PipelineService) GetStatus(ctx context.Context, stage PipelineStage) (string, error) {
	var status string
	err := ps.db.NewRaw(`SELECT status FROM pipeline_stages WHERE stage = ?`, stage).Scan(ctx, &status)
	return status, err
}

func (ps *PipelineService) GetAllStatuses(ctx context.Context) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := ps.db.NewRaw(`SELECT stage, status, started_at, completed_at, error_message, metadata FROM pipeline_stages ORDER BY CASE stage WHEN 'migrations' THEN 1 WHEN 'jsonl_load' THEN 2 WHEN 'embeddings_load' THEN 3 WHEN 'graph_vertices' THEN 4 WHEN 'graph_edges' THEN 5 END`).Scan(ctx, &results)
	return results, err
}

func (ps *PipelineService) ResetStage(ctx context.Context, stage PipelineStage) error {
	_, err := ps.db.ExecContext(ctx, `UPDATE pipeline_stages SET status = 'pending', started_at = NULL, completed_at = NULL, error_message = NULL, metadata = '{}' WHERE stage = ?`, stage)
	return err
}
