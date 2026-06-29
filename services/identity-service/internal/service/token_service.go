package service

import (
	"fmt"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/identity-service/internal/domain"
)

type Claims struct {
	UserID   string          `json:"user_id"`
	Email    string          `json:"email"`
	Role     domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type TokenService interface {
	GenerateAccessToken(user *domain.User) (string, error)
	GenerateRefreshToken() (string, int, error)
	ValidateAccessToken(token string) (*Claims, error)
	HashToken(token string) string
}

type tokenService struct {
	secret     string
	expiry     time.Duration
}

func NewTokenService(secret string, expiry time.Duration) TokenService {
	return &tokenService{
		secret: secret,
		expiry: expiry,
	}
}

func (s *tokenService) GenerateAccessToken(user *domain.User) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "spark-identity",
			Subject:   user.ID.String(),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (s *tokenService) GenerateRefreshToken() (string, int, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	refreshToken := hex.EncodeToString(b)
	return refreshToken, len(refreshToken), nil
}

func (s *tokenService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

func (s *tokenService) HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
