package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	jwtlib "github.com/golang-jwt/jwt/v5"

	"auth/jwt"
)

type AuthMiddleware struct {
	generator *jwt.Generator
}

func NewAuthMiddleware(generator *jwt.Generator) *AuthMiddleware {
	return &AuthMiddleware{generator: generator}
}

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := m.extractToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := m.validateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := m.addToContext(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func (m *AuthMiddleware) validateToken(tokenString string) (*jwt.Claims, error) {
	claims := &jwt.Claims{}

	token, err := jwtlib.ParseWithClaims(tokenString, claims,
		func(token *jwtlib.Token) (interface{}, error) {
			return m.generator.PublicKey, nil
		})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}

func (m *AuthMiddleware) addToContext(ctx context.Context, claims *jwt.Claims) context.Context {
	ctx = context.WithValue(ctx, "user_id", claims.UserID)
	ctx = context.WithValue(ctx, "user_email", claims.Email)
	return ctx
}
