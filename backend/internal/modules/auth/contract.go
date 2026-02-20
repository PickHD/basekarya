package auth

import (
	"context"
	"basekarya-backend/internal/modules/user"
)

type Hasher interface {
	CheckPasswordHash(password, hash string) bool
}

type TokenProvider interface {
	GenerateToken(userID uint, role string, employeeID *uint) (string, error)
}

type UserProvider interface {
	FindByUsername(ctx context.Context, username string) (*user.User, error)
}
