package domain

import (
	"fmt"
	"net/mail"
	"strings"
)

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email")
	}

	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format")
	}

	domain := parts[1]
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("invalid email: domain must contain a dot (e.g., example.com)")
	}

	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return fmt.Errorf("invalid email: domain cannot start or end with dot")
	}

	if parts[0] == "" {
		return fmt.Errorf("invalid email: local part cannot be empty")
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if len(password) < 7 {
		return fmt.Errorf("password cannot be shorter than 7 characters")
	}
	return nil
}
