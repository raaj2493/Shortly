package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raaj2493/Shortly/Backend/config"
)

func Connect(ctx context.Context , cfg *config.Config ) (*pgxpool.Pool , error){
	// Build the connection string (DSN).
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	// Parse the DSN into a pgxpool configuration.
	poolConfig , err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	// Optional: tune the pool (defaults are fine for now)
	poolConfig.MaxConns = 10                     // maximum number of connections
	poolConfig.MinConns = 2                      // minimum idle connections
	poolConfig.MaxConnLifetime = 1 * time.Hour   // how long a connection can live
	poolConfig.MaxConnIdleTime = 30 * time.Minute // how long an idle connection is kept

	// Create the pool.
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Ping the database to verify the connection.
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close() // clean up if ping fails
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}

