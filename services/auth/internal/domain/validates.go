package domain

import (
	"fmt"
	"net/mail"
)

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email")
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
