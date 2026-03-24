package auth

import (
	"basekarya-backend/internal/modules/user"
	"context"
)

type Hasher interface {
	CheckPasswordHash(password, hash string) bool
}

type TokenProvider interface {
	GenerateToken(userID uint, role string, employeeID *uint, permissions []string) (string, error)
}

type UserProvider interface {
	FindByUsername(ctx context.Context, username string) (*user.User, error)
}
