package infrastructure

import (
	"errors"
	"fmt"
	"hris-backend/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtProvider struct {
	secretKey      string
	issuer         string
	expireDuration time.Duration
}

func NewJWTProvider(cfg *config.Config) *JwtProvider {
	return &JwtProvider{
		secretKey:      cfg.JWT.Secret,
		issuer:         "hris-app",
		expireDuration: time.Hour * time.Duration(cfg.JWT.ExpiresIn),
	}
}

type MyClaims struct {
	UserID     uint   `json:"user_id"`
	Role       string `json:"role"`
	EmployeeID *uint  `json:"employee_id"`
	jwt.RegisteredClaims
}

func (p *JwtProvider) GenerateToken(userID uint, role string, employeeID *uint) (string, error) {
	claims := &MyClaims{
		UserID:     userID,
		Role:       role,
		EmployeeID: employeeID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.expireDuration)),
			Issuer:    p.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.secretKey))
}

func (p *JwtProvider) ValidateToken(tokenString string) (*MyClaims, error) {
	claims := &MyClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(p.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
