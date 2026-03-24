package di

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/uptrace/bun/migrate"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations executes all pending UP database migrations.
func RunMigrations(lc fx.Lifecycle, db *Database, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Running database migrations...")

			// Get the migrations subdirectory from embedded FS
			subFS, err := fs.Sub(migrationsFS, "migrations")
			if err != nil {
				return err
			}

			// Discover embedded SQL migrations
			migrations := migrate.NewMigrations()
			if err := migrations.Discover(subFS); err != nil {
				return err
			}

			// Create migrator with custom table name
			migrator := migrate.NewMigrator(db.DB, migrations,
				migrate.WithTableName("migrations"),
			)

			// Create migrations tracking table
			if err := migrator.Init(ctx); err != nil {
				return err
			}

			// Run UP migrations only
			group, err := migrator.Migrate(ctx)
			if err != nil {
				return err
			}

			if group.ID == 0 {
				logger.Info("No new migrations to apply")
			} else {
				logger.Info("Migrations completed successfully",
					zap.String("group", fmt.Sprintf("%d", group.ID)))
			}

			return nil
		},
	})
}
