package database

import (
	"context"

	"github.com/Masterminds/squirrel"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	PgxBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type Database struct {
	Pool    *pgxpool.Pool
	Logger  *zap.Logger
	Builder *squirrel.StatementBuilderType
}

func NewPgxPool(connUrl string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connUrl)
	if err != nil {
		return nil, err
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, err
}

func New(pool *pgxpool.Pool, builder *squirrel.StatementBuilderType, logger *zap.Logger) (*Database, error) {
	logger.Info("Database started")
	logger.Debug(pool.Config().ConnString())

	return &Database{ 
		Pool:    pool,
		Logger:  logger,
		Builder: builder,
	}, nil
}

func (db *Database) Close() {
	db.Pool.Close()
}

func NewPostgres(connUrl string, logger *zap.Logger) (*Database, error) {
	pool, err := NewPgxPool(connUrl)
	if err != nil {
		return nil, err
	}

	return New(pool, &PgxBuilder, logger)
}
