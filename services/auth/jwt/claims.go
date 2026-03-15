package jwt

import (
	"auth/internal/config"
	"auth/internal/domain"
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewClaimsFromUser(user *domain.User) *Claims {
	return &Claims{
		UserID: user.ID,
		Email:  user.Email,
	}
}

func (c *Claims) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "user_id", c.UserID)
}

type Generator struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	config     *config.JWTConfig
}

func NewGenerator(cfg *config.JWTConfig) (*Generator, error) {

	privateKey, publicKey, err := loadKeys(cfg.PrivateKeyPath, cfg.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load keys: %w", err)
	}

	return &Generator{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		config:     cfg,
	}, nil
}

func (g *Generator) GenerateAccessToken(user *domain.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.config.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    g.config.Issuer,
			Subject:   "access",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(g.PrivateKey)
}

func (g *Generator) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return g.PublicKey, nil
		})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("refresh token is not valid")
	}

	// Проверяем что это refresh token
	if claims.Subject != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	return claims, nil
}

func (g *Generator) GenerateRefreshToken(user *domain.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.config.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    g.config.Issuer,
			Subject:   "refresh",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(g.PrivateKey)
}
