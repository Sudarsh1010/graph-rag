package di

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/sudarsh1010/graph-rag/internal/graph"
	"github.com/sudarsh1010/graph-rag/internal/loader"
	"github.com/sudarsh1010/graph-rag/internal/store"
	"github.com/uptrace/bun/migrate"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func RunPipeline(lc fx.Lifecycle, cfg *Config, db *Database, gs *graph.Service, ps *PipelineService, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			start := time.Now()
			logger.Info("starting pipeline...")

			if err := ps.Init(ctx); err != nil {
				return fmt.Errorf("pipeline init: %w", err)
			}

			// Stage 1: Migrations
			if err := ps.RunStage(ctx, StageMigrations, func(ctx context.Context) error {
				return runMigrations(ctx, db, logger)
			}); err != nil {
				return err
			}

			// Stage 2: JSONL Load
			if err := ps.RunStage(ctx, StageJSONLLoad, func(ctx context.Context) error {
				coord := loader.NewCoordinator()
				result, err := coord.LoadAll(ctx, loader.LoadConfig{
					DatasetPath: cfg.DatasetPath,
					DB:          db.DB,
					Logger:      logger,
				})
				if err != nil {
					return err
				}
				logger.Info("data loaded",
					zap.Int("total", result.Total),
					zap.Int("inserted", result.Inserted),
					zap.Int("skipped", result.Skipped),
					zap.Int("errors", len(result.Errors)),
				)
				return nil
			}); err != nil {
				return err
			}

			// Stage 3: Embeddings Load (skip if directory doesn't exist)
			if _, err := os.Stat(cfg.EmbeddingsPath); err == nil {
				if err := ps.RunStage(ctx, StageEmbeddingsLoad, func(ctx context.Context) error {
					return store.LoadEmbeddings(ctx, db.DB, cfg.EmbeddingsPath)
				}); err != nil {
					return err
				}
			} else {
				logger.Warn("embeddings directory not found, skipping",
					zap.String("path", cfg.EmbeddingsPath),
				)
			}

			// Stage 4: Graph Vertices
			if err := ps.RunStage(ctx, StageGraphVertices, func(ctx context.Context) error {
				if err := gs.Initialize(ctx); err != nil {
					return fmt.Errorf("graph init: %w", err)
				}
				return gs.LoadVertices(ctx)
			}); err != nil {
				return err
			}

			// Stage 5: Graph Edges
			if err := ps.RunStage(ctx, StageGraphEdges, func(ctx context.Context) error {
				return gs.LoadEdges(ctx)
			}); err != nil {
				return err
			}

			logger.Info("pipeline completed",
				zap.Duration("total_elapsed", time.Since(start)),
			)
			return nil
		},
	})
}

func runMigrations(ctx context.Context, db *Database, logger *zap.Logger) error {
	subFS, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	migrations := migrate.NewMigrations()
	if err := migrations.Discover(subFS); err != nil {
		return err
	}

	migrator := migrate.NewMigrator(db.DB, migrations, migrate.WithTableName("migrations"))
	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return err
	}

	if group.ID == 0 {
		logger.Info("No new migrations to apply")
	} else {
		logger.Info("Migrations completed successfully",
			zap.String("group", fmt.Sprintf("%d", group.ID)),
		)
	}

	return nil
}
