package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthentication interface {
	GenerateToken(userID int64) (string, error)
	GetUserID(token string) (int64, error)
	GeneratePasswordHash(password string) (string, error)
	CheckPasswordHash(password, hash string) error
}

const (
	tokenTTL    = 1 * time.Hour
	hashingCost = 10
)

type auth struct {
	secretKey string
	tokenTTL  time.Duration
}

func New(secretKey string) *auth {
	return &auth{
		secretKey: secretKey,
		tokenTTL:  tokenTTL,
	}
}

type Claims struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"`
}

func (a *auth) GenerateToken(userID int64) (string, error) {
	claims := Claims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken, nil
}

func (a *auth) GetUserID(tokenString string) (int64, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(a.secretKey), nil
		})
	if err != nil {
		return 0, fmt.Errorf("failed to parse token with claims, %w", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("token not valid")
	}

	return claims.UserID, nil
}

func (a *auth) GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), hashingCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate from password, %w", err)
	}
	return string(bytes), nil
}

func (a *auth) CheckPasswordHash(password, hash string) error {
	//nolint // Не за чем оборачивать ошибку
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
