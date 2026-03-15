package repository

import (
	"auth/internal/config"
	"auth/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupDatabase(ctx context.Context, cfg *config.DBConfig) (*pgxpool.Pool, *PostgresUserRepo, *domain.Service) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		panic(err)
	}

	poolConfig.MaxConns = cfg.MaxConnections
	poolConfig.MaxConnIdleTime = cfg.IdleTimeout

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		panic(err)
	}

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	fmt.Println("connected to database successfully")

	userRepo := NewPostgresUserRepo(pool)
	userService := domain.NewService(userRepo)

	return pool, userRepo, userService
}
