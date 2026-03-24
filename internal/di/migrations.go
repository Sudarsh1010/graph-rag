package di

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// RunMigrations executes all pending database migrations.
func RunMigrations(lc fx.Lifecycle, db *Database, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Running database migrations...")

			// Get migrations directory path
			migrationsDir := "internal/infrastructure/persistence/bun/migrations"

			// Read migration files
			entries, err := os.ReadDir(migrationsDir)
			if err != nil {
				logger.Warn(
					"Migrations directory not found, skipping migrations",
					zap.String("path", migrationsDir), zap.Error(err),
				)
				return nil
			}

			// Execute each migration file
			for _, entry := range entries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
					continue
				}

				filePath := filepath.Join(migrationsDir, entry.Name())
				content, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf(
						"failed to read migration file %s: %w",
						entry.Name(),
						err,
					)
				}

				logger.Info(
					"Executing migration",
					zap.String("file", entry.Name()),
				)

				_, err = db.ExecContext(ctx, string(content))
				if err != nil {
					return fmt.Errorf(
						"failed to execute migration %s: %w",
						entry.Name(), err,
					)
				}
			}

			logger.Info("Database migrations completed successfully")
			return nil
		},
	})
}
