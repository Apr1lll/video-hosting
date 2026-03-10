package repository

import (
	"auth/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupDatabase(ctx context.Context) (*pgxpool.Pool, *PostgresUserRepo, *domain.Service) {
	connStr := "postgres://postgres:postgres@localhost:5432?sslmode=disable"

	pool, err := pgxpool.New(ctx, connStr)
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
