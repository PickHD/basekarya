package user

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"
)

type Service interface {
	GetProfile(userID uint) (*UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest, file *multipart.FileHeader) error
	ChangePassword(userID uint, req *ChangePasswordRequest) error
}

type service struct {
	repo    Repository
	bcrypt  Hasher
	storage StorageProvider
}

func NewService(repo Repository, bcrypt Hasher, storage StorageProvider) Service {
	return &service{repo, bcrypt, storage}
}

func (s *service) GetProfile(userID uint) (*UserProfileResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	resp := &UserProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}

	if user.Employee != nil {
		resp.FullName = user.Employee.FullName
		resp.NIK = user.Employee.NIK
		resp.PhoneNumber = user.Employee.PhoneNumber
		resp.ProfilePictureUrl = user.Employee.ProfilePictureUrl

		if user.Employee.Department != nil {
			resp.DepartmentName = user.Employee.Department.Name
		}
		if user.Employee.Shift != nil {
			resp.ShiftName = user.Employee.Shift.Name
		}
	} else {
		resp.FullName = "Super Administrator"
	}

	return resp, nil
}

func (s *service) UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest, file *multipart.FileHeader) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	if user.Employee == nil {
		return errors.New("employee data not found")
	}

	user.Employee.PhoneNumber = req.PhoneNumber

	if file != nil {
		fileName := fmt.Sprintf("users/%d/profile-%d.jpg", userID, time.Now().Unix())
		fileURL, err := s.storage.UploadFileMultipart(ctx, file, fileName)
		if err != nil {
			return err
		}

		user.Employee.ProfilePictureUrl = fileURL
	}

	return s.repo.UpdateEmployee(user.Employee)
}

func (s *service) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	// check hash if user trying to change password again
	if !user.MustChangePassword {
		if !s.bcrypt.CheckPasswordHash(req.OldPassword, user.PasswordHash, false) {
			return errors.New("invalid old password")
		}
	}

	hashedPassword, err := s.bcrypt.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	user.MustChangePassword = false

	return s.repo.UpdateUser(user)
}
