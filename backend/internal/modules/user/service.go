package user

import (
	"context"
	"errors"
	"fmt"
	"hris-backend/pkg/constants"
	"hris-backend/pkg/response"
	"mime/multipart"
	"time"
)

type Service interface {
	GetProfile(userID uint) (*UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest, file *multipart.FileHeader) error
	ChangePassword(userID uint, req *ChangePasswordRequest) error
	GetAllEmployees(ctx context.Context, page, limit int, search string) ([]EmployeeListResponse, *response.Meta, error)
	CreateEmployee(ctx context.Context, req *CreateEmployeeRequest) error
	UpdateEmployee(ctx context.Context, id uint, req *UpdateEmployeeRequest) error
	DeleteEmployee(ctx context.Context, id uint) error
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
			resp.ShiftStartTime = user.Employee.Shift.StartTime
			resp.ShiftEndTime = user.Employee.Shift.EndTime
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

	if req.PhoneNumber != "" {
		user.Employee.PhoneNumber = req.PhoneNumber
	}

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

	if !s.bcrypt.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		return errors.New("invalid old password")
	}

	hashedPassword, err := s.bcrypt.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	user.MustChangePassword = false

	return s.repo.UpdateUser(user)
}

func (s *service) GetAllEmployees(ctx context.Context, page, limit int, search string) ([]EmployeeListResponse, *response.Meta, error) {
	users, total, err := s.repo.FindAllEmployees(page, limit, search)
	if err != nil {
		return nil, nil, err
	}

	if len(users) == 0 {
		return []EmployeeListResponse{}, nil, nil
	}

	var list []EmployeeListResponse
	for _, u := range users {
		deptName := "-"
		shiftName := "-"
		if u.Employee != nil {
			if u.Employee.Department != nil {
				deptName = u.Employee.Department.Name
			}
			if u.Employee.Shift != nil {
				shiftName = u.Employee.Shift.Name
			}

			list = append(list, EmployeeListResponse{
				ID:             u.Employee.ID,
				FullName:       u.Employee.FullName,
				NIK:            u.Employee.NIK,
				Username:       u.Username,
				DepartmentName: deptName,
				ShiftName:      shiftName,
			})
		}
	}

	meta := response.NewMeta(page, limit, total)
	return list, meta, nil
}

func (s *service) CreateEmployee(ctx context.Context, req *CreateEmployeeRequest) error {
	tx := s.repo.StartTX()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	checkUser, err := s.repo.FindByUsername(req.Username)
	if err == nil && checkUser.ID != 0 {
		return errors.New("username already exists")
	}

	hashPass, _ := s.bcrypt.HashPassword(req.Username)

	newUser := User{
		Username:           req.Username,
		PasswordHash:       hashPass,
		Role:               string(constants.UserRoleEmployee),
		MustChangePassword: true,
	}

	if err := s.repo.CreateUser(tx, &newUser); err != nil {
		return err
	}

	newEmp := Employee{
		UserID:       newUser.ID,
		FullName:     req.FullName,
		NIK:          req.NIK,
		DepartmentID: req.DepartmentID,
		ShiftID:      req.ShiftID,
	}

	if err := s.repo.CreateEmployee(tx, &newEmp); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateEmployee(ctx context.Context, id uint, req *UpdateEmployeeRequest) error {
	emp, err := s.repo.FindEmployeeByID(id)
	if err != nil {
		return errors.New("employee not found")
	}

	if req.FullName != "" {
		emp.FullName = req.FullName
	}
	if req.NIK != "" {
		emp.NIK = req.NIK
	}
	if req.DepartmentID > 0 {
		emp.DepartmentID = req.DepartmentID
	}
	if req.ShiftID > 0 {
		emp.ShiftID = req.ShiftID
	}

	return s.repo.UpdateEmployee(emp)
}

func (s *service) DeleteEmployee(ctx context.Context, id uint) error {
	emp, err := s.repo.FindEmployeeByID(id)
	if err != nil {
		return errors.New("employee not found")
	}

	return s.repo.DeleteUser(emp.UserID)
}
