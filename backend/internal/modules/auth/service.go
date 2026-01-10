package auth

import (
	"errors"
	"hris-backend/internal/modules/user"
)

type Service interface {
	Login(username, password string) (*LoginResponse, error)
}

type service struct {
	userRepo      user.Repository
	hasher        Hasher
	tokenProvider TokenProvider
}

func NewService(userRepo user.Repository, hasher Hasher, tokenProvider TokenProvider) Service {
	return &service{
		userRepo:      userRepo,
		hasher:        hasher,
		tokenProvider: tokenProvider,
	}
}

func (s *service) Login(username, password string) (*LoginResponse, error) {
	foundUser, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !s.hasher.CheckPasswordHash(password, foundUser.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	tokenString, err := s.tokenProvider.GenerateToken(foundUser.ID, foundUser.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:              tokenString,
		MustChangePassword: foundUser.MustChangePassword,
	}, nil
}
