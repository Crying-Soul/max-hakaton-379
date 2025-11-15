package repository

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	dbsqlc "maxBot/internal/db/sqlc"
)

// Repository manages database access and migrations for the application.
type Repository struct {
	pool    *pgxpool.Pool
	queries *dbsqlc.Queries
}

// New creates a repository by connecting to the database and applying migrations.
func New(ctx context.Context, dsn, migrationsPath string) (*Repository, error) {
	if dsn == "" {
		return nil, fmt.Errorf("database dsn is required")
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if migrationsPath == "" {
		migrationsPath = "db/migrations"
	}

	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dsn)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("init migrate: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		pool.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return &Repository{
		pool:    pool,
		queries: dbsqlc.New(pool),
	}, nil
}

// Queries exposes the generated sqlc querier for downstream services.
func (r *Repository) Queries() dbsqlc.Querier {
	return r.queries
}

// Close releases the underlying database connection pool.
func (r *Repository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}
