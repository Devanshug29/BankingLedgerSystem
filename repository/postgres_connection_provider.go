package repository

import (
	"BankingLedgerSystem/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectionProviderInterface defines the contract for DB connections
type ConnectionProviderInterface interface {
	GetConnection(ctx context.Context) *pgxpool.Pool
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Close() error
}

// dbClient wraps pgxpool for connection tracking
type dbClient struct {
	db *pgxpool.Pool
}

// databaseMap maintains active DB connections and their stats
type databaseMap struct {
	connections map[string]*dbClient
}

var dbMap = &databaseMap{
	connections: make(map[string]*dbClient),
}

// DBConfig holds Postgres connection settings
type DBConfig struct {
	URL                       string
	MinConnections            int32
	MaxConnections            int32
	MaxConnectionIdleTimeInMS int32
	ConnectionTimeoutInMS     int32
}

// ConnectionProvider implements ConnectionProviderInterface
type ConnectionProvider struct {
	config *config.Config
	db     *pgxpool.Pool
}

// InitDB initializes and verifies the Postgres connection pool
func (c *ConnectionProvider) InitDB(ctx context.Context, name string) error {
	poolConfig, err := c.getPoolConfig()
	if err != nil {
		return err
	}
	c.db, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("failed to create pgx pool: %w", err)
	}

	if err = c.db.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	db := &dbClient{db: c.db}
	dbMap.connections[name] = db
	return nil
}

// GetConnection returns the pgxpool connection
func (c *ConnectionProvider) GetConnection(ctx context.Context) *pgxpool.Pool {
	return c.db
}

// BeginTx starts a new transaction with options
func (c *ConnectionProvider) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return c.db.BeginTx(ctx, opts)
}

// Begin starts a simple transaction
func (c *ConnectionProvider) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.db.Begin(ctx)
}

// Close shuts down the connection pool
func (c *ConnectionProvider) Close() error {
	if c.db != nil {
		c.db.Close()
	}
	return nil
}

// getPoolConfig builds pgxpool config from DBConfig
func (c *ConnectionProvider) getPoolConfig() (*pgxpool.Config, error) {
	config, err := pgxpool.ParseConfig(c.config.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pg db config connection url: %w", err)
	}

	return config, nil
}

// NewConnectionProvider creates a new provider instance
func NewConnectionProvider(cfg *config.Config) *ConnectionProvider {
	return &ConnectionProvider{config: cfg}
}
