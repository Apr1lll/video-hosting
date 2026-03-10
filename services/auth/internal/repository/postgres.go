package repository

import (
	"auth/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepo(pool *pgxpool.Pool) *PostgresUserRepo {
	return &PostgresUserRepo{pool: pool}
}

func (r *PostgresUserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO users(email, password_hash, created_at)
	VALUES($1, $2, $3)
	RETURNING id
	`

	err := r.pool.QueryRow(ctx, query, user.Email, user.PasswordHash, user.CreatedAt).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
	SELECT id, email, password_hash, created_at
	FROM users
	WHERE email = $1
	`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
	SELECT id, email, password_hash, created_at
	FROM users
	WHERE id = $1
	`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}
