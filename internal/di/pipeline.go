package di

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	db     *Database
	logger *zap.Logger
}

func NewPipelineService(db *Database, logger *zap.Logger) *PipelineService {
	return &PipelineService{db: db, logger: logger}
}

func (ps *PipelineService) Init(ctx context.Context) error {
	_, err := ps.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS pipeline_stages (
			stage VARCHAR(50) PRIMARY KEY,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			started_at TIMESTAMPTZ,
			completed_at TIMESTAMPTZ,
			error_message TEXT,
			metadata JSONB DEFAULT '{}'
		)`)
	return err
}

func (ps *PipelineService) RunStage(ctx context.Context, stage PipelineStage, fn func(ctx context.Context) error) error {
	stageName := string(stage)

	var status string
	err := ps.db.DB.NewRaw("SELECT status FROM pipeline_stages WHERE stage = ?", stageName).Scan(ctx, &status)

	if err == sql.ErrNoRows {
		_, err = ps.db.ExecContext(ctx, "INSERT INTO pipeline_stages (stage, status) VALUES (?, 'pending')", stageName)
		if err != nil {
			return fmt.Errorf("insert pipeline stage %s: %w", stageName, err)
		}
		status = "pending"
	} else if err != nil {
		return fmt.Errorf("query pipeline stage %s: %w", stageName, err)
	}

	if status == "completed" {
		ps.logger.Info("pipeline: skipping stage (already completed)", zap.String("stage", stageName))
		return nil
	}

	if status == "failed" || status == "running" {
		ps.logger.Info("pipeline: retrying stage", zap.String("stage", stageName), zap.String("previous_status", status))
	}

	_, _ = ps.db.ExecContext(ctx,
		"UPDATE pipeline_stages SET status = 'running', started_at = NOW(), error_message = NULL WHERE stage = ?",
		stageName,
	)

	start := time.Now()
	runErr := fn(ctx)
	elapsed := time.Since(start)

	if runErr != nil {
		ps.db.ExecContext(ctx,
			"UPDATE pipeline_stages SET status = 'failed', error_message = ? WHERE stage = ?",
			runErr.Error(), stageName,
		)
		return fmt.Errorf("stage %s failed: %w", stageName, runErr)
	}

	ps.db.ExecContext(ctx,
		"UPDATE pipeline_stages SET status = 'completed', completed_at = NOW(), error_message = NULL WHERE stage = ?",
		stageName,
	)

	ps.logger.Info("pipeline: completed stage", zap.String("stage", stageName), zap.Duration("elapsed", elapsed))
	return nil
}

func (ps *PipelineService) GetStatus(ctx context.Context, stage PipelineStage) (string, error) {
	var status string
	err := ps.db.DB.QueryRowContext(ctx,
		"SELECT status FROM pipeline_stages WHERE stage = ?", string(stage),
	).Scan(&status)
	if err == sql.ErrNoRows {
		return "not_found", nil
	}
	if err != nil {
		return "", fmt.Errorf("get pipeline stage status: %w", err)
	}
	return status, nil
}

func (ps *PipelineService) ResetStage(ctx context.Context, stage PipelineStage) error {
	_, err := ps.db.ExecContext(ctx,
		"UPDATE pipeline_stages SET status = 'pending', started_at = NULL, completed_at = NULL, error_message = NULL, metadata = '{}' WHERE stage = ?",
		stage,
	)
	return err
}

func (ps *PipelineService) GetAllStatuses(ctx context.Context) ([]map[string]interface{}, error) {
	var rows []map[string]interface{}
	err := ps.db.NewRaw("SELECT stage, status, started_at, completed_at, error_message, metadata FROM pipeline_stages ORDER BY CASE stage WHEN 'migrations' THEN 1 WHEN 'jsonl_load' THEN 2 WHEN 'embeddings_load' THEN 3 WHEN 'graph_vertices' THEN 4 WHEN 'graph_edges' THEN 5 END").Scan(ctx, &rows)
	return rows, err
}
