package auth

import (
	"hris-backend/internal/modules/user"
)

type Hasher interface {
	CheckPasswordHash(password, hash string) bool
}

type TokenProvider interface {
	GenerateToken(userID uint, role string, employeeID *uint) (string, error)
}

type UserProvider interface {
	FindByUsername(username string) (*user.User, error)
}
