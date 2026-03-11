package jwt

import (
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
	Issuer     string
}

func NewGenerator(privatePath, publicPath string) (*Generator, error) {
	privateKey, publicKey, err := loadKeys(privatePath, publicPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load keys: %w", err)
	}

	return &Generator{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Issuer:     "auth-service",
	}, nil
}

func (g *Generator) GenerateAccessToken(user *domain.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    g.Issuer,
			Subject:   "access",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(g.PrivateKey)
}

func (g *Generator) GenerateRefreshToken(user *domain.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    g.Issuer,
			Subject:   "refresh",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(g.PrivateKey)
}
