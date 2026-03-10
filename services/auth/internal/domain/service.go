package domain

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) (*User, error) {

	if err := ValidateEmail(email); err != nil {
		return nil, err
	}

	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}
	if existingUser != nil {
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

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user")
	}

	return user, nil

}

func (s *Service) Login(ctx context.Context, email, password string) (*User, error) {
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}

	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid email or password")
	}

	return user, nil
}
