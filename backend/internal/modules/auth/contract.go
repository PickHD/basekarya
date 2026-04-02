package auth

import (
	"basekarya-backend/internal/modules/user"
	"context"
	"time"
)

type Hasher interface {
	CheckPasswordHash(password, hash string) bool
	HashPassword(password string) (string, error)
}

type TokenProvider interface {
	GenerateToken(userID uint, role string, employeeID *uint, permissions []string) (string, error)
}

type UserProvider interface {
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	FindEmployeeByEmail(ctx context.Context, email string) (*user.Employee, error)
	UpdatePasswordByEmail(ctx context.Context, email string, password string) error
}

type CacheProvider interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}

type EmailProvider interface {
	Send(to, subject, htmlBody string) error
}
