package auth

import (
	"context"
	"errors"
	"basekarya-backend/pkg/constants"
)

type Service interface {
	Login(username, password string) (*LoginResponse, error)
}

type service struct {
	user          UserProvider
	hasher        Hasher
	tokenProvider TokenProvider
}

func NewService(user UserProvider, hasher Hasher, tokenProvider TokenProvider) Service {
	return &service{
		user:          user,
		hasher:        hasher,
		tokenProvider: tokenProvider,
	}
}

func (s *service) Login(username, password string) (*LoginResponse, error) {
	foundUser, err := s.user.FindByUsername(context.Background(), username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !s.hasher.CheckPasswordHash(password, foundUser.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	var employeeID *uint
	if foundUser.Employee != nil && foundUser.Role != nil && foundUser.Role.Name != string(constants.UserRoleSuperadmin) {
		employeeID = &foundUser.Employee.ID
	}

	var permissions []string
	if foundUser.Role != nil {
		for _, permission := range foundUser.Role.Permissions {
			permissions = append(permissions, permission.Name)
		}
	}

	tokenString, err := s.tokenProvider.GenerateToken(foundUser.ID, foundUser.Role.Name, employeeID, permissions)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:              tokenString,
		MustChangePassword: foundUser.MustChangePassword,
	}, nil
}
