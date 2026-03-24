package di

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver registration
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.uber.org/zap"
)

type Database struct {
	*bun.DB
}

func NewDatabase(
	cfg *Config,
	logger *zap.Logger,
) (*Database, error) {
	ctx := context.Background()

	// Build DSN from config - supports either full DSN or individual params
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUsername,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresSSLMode,
	)

	logger.Info("Initializing PostgreSQL database", zap.String("host", cfg.PostgresHost))

	// Open connection using pgx driver
	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	// Apply connection pooling
	sqldb.SetMaxOpenConns(BunMaxOpenConns)
	sqldb.SetMaxIdleConns(BunMaxIdleConns)
	sqldb.SetConnMaxLifetime(BunConnMaxLifetime)
	sqldb.SetConnMaxIdleTime(BunConnMaxIdleTime)

	// Create Bun DB instance with PostgreSQL dialect
	db := bun.NewDB(sqldb, pgdialect.New())

	// Test connection
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database initialized successfully")

	return &Database{db}, nil
}

func (d *Database) Close(_ context.Context) error {
	return d.DB.Close()
}
