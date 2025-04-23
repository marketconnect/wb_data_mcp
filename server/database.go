package server

import (
	"context"
	"fmt"
	"time"

	"github.com/marketconnect/wb_data_mcp/config"

	"github.com/marketconnect/db_client/clickhouse"
	"github.com/marketconnect/db_client/postgresql"
)

// Database holds all database connections
type Database struct {
	ClickHouse clickhouse.ClickHouseClient
	PostgreSQL postgresql.PostgreSQLClient
}

// NewDatabase initializes all database connections
func NewDatabase(ctx context.Context, cfg *config.Config) (*Database, error) {
	// Initialize ClickHouse client
	chConfig := clickhouse.NewClickHouseConfig(
		cfg.Clickhouse.Host,
		cfg.Clickhouse.Port,
		cfg.Clickhouse.Database,
		cfg.Clickhouse.Username,
		cfg.Clickhouse.Password,
	)
	chClient, err := clickhouse.NewClickHouseClient(ctx, 5, time.Second*5, chConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ClickHouse client: %w", err)
	}

	// Initialize PostgreSQL client
	pgConfig := postgresql.NewPgConfig(
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)
	pgClient, err := postgresql.NewClient(ctx, 5, time.Second*5, pgConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL client: %w", err)
	}

	return &Database{
		ClickHouse: chClient,
		PostgreSQL: pgClient,
	}, nil
}

// Close closes all database connections
func (db *Database) Close() error {
	var errs []error

	if err := db.ClickHouse.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close ClickHouse connection: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing database connections: %v", errs)
	}
	return nil
}
