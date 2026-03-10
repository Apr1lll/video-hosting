package repository

import (
	"auth/internal/domain"
	"context"
	"database/sql"
	"fmt"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) Register(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO users(email, password_hash, created_at)
	VALUES($1, $2, $3)
	`

	err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
	SELECT id, email, user, password_hash, created_at
	FROM users
	WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not exists")
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
	err := r.db.QueryRowContext(ctx, query, user.Email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not exists")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID")
	}

	return user, nil
}
