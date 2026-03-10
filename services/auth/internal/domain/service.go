package domain

import (
	"auth/internal/repository"
	"context"
	"fmt"
	"time"
)

type UserActions interface {
	Register(ctx context.Context, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, error)
}

type Service struct {
	act  UserActions
	repo repository.PostgresUserRepo
}

func (s *Service) Register(ctx context.Context, email, password string) (*User, error) {

	if err := ValidateEmail(email); err != nil {
		return nil, err
	}

	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	_, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hashing password: %w", err)
	}

	user := &User{
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.Register(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user")
	}

	return user, nil

}
