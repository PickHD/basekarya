package user

import (
	"context"
	"errors"
	"fmt"
	"hris-backend/internal/infrastructure"
	"hris-backend/pkg/constants"
	"hris-backend/pkg/response"
	"mime/multipart"
	"time"
)

type Service interface {
	GetProfile(userID uint) (*UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest, file *multipart.FileHeader) error
	ChangePassword(ctx context.Context, userID uint, req *ChangePasswordRequest) error
	GetAllEmployees(ctx context.Context, page, limit int, search string) ([]EmployeeListResponse, *response.Meta, error)
	CreateEmployee(ctx context.Context, req *CreateEmployeeRequest) error
	UpdateEmployee(ctx context.Context, id uint, req *UpdateEmployeeRequest) error
	DeleteEmployee(ctx context.Context, id uint) error
}

type service struct {
	repo               Repository
	bcrypt             Hasher
	storage            StorageProvider
	leaveGenerator     LeaveBalanceGenerator
	transactionManager infrastructure.TransactionManager
}

func NewService(repo Repository, bcrypt Hasher, storage StorageProvider, leaveGenerator LeaveBalanceGenerator, transactionManager infrastructure.TransactionManager) Service {
	return &service{repo, bcrypt, storage, leaveGenerator, transactionManager}
}

func (s *service) GetProfile(userID uint) (*UserProfileResponse, error) {
	user, err := s.repo.FindByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	resp := &UserProfileResponse{
		ID:                 user.ID,
		Username:           user.Username,
		Role:               user.Role,
		MustChangePassword: user.MustChangePassword,
	}

	if user.Employee != nil {
		resp.FullName = user.Employee.FullName
		resp.NIK = user.Employee.NIK
		resp.PhoneNumber = user.Employee.PhoneNumber
		resp.ProfilePictureUrl = user.Employee.ProfilePictureUrl
		resp.BankName = user.Employee.BankName
		resp.BaseSalary = user.Employee.BaseSalary
		resp.BankAccountNumber = user.Employee.BankAccountNumber
		resp.BankAccountHolder = user.Employee.BankAccountHolder
		resp.NPWP = user.Employee.NPWP

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
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Employee == nil {
		return errors.New("employee data not found")
	}

	user, err = s.buildEmployeeData(ctx, user, req, file)
	if err != nil {
		return err
	}

	return s.repo.UpdateEmployee(ctx, user.Employee)
}

func (s *service) ChangePassword(ctx context.Context, userID uint, req *ChangePasswordRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
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

	return s.repo.UpdateUser(ctx, user)
}

func (s *service) GetAllEmployees(ctx context.Context, page, limit int, search string) ([]EmployeeListResponse, *response.Meta, error) {
	users, total, err := s.repo.FindAllEmployees(ctx, page, limit, search)
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
		baseSalary := 0.0

		if u.Employee != nil {
			if u.Employee.Department != nil {
				deptName = u.Employee.Department.Name
			}
			if u.Employee.Shift != nil {
				shiftName = u.Employee.Shift.Name
			}
			if u.Employee.BaseSalary != 0 {
				baseSalary = u.Employee.BaseSalary
			}

			list = append(list, EmployeeListResponse{
				ID:             u.Employee.ID,
				FullName:       u.Employee.FullName,
				NIK:            u.Employee.NIK,
				Username:       u.Username,
				DepartmentName: deptName,
				ShiftName:      shiftName,
				BaseSalary:     baseSalary,
			})
		}
	}

	meta := response.NewMetaOffset(page, limit, total)
	return list, meta, nil
}

func (s *service) CreateEmployee(ctx context.Context, req *CreateEmployeeRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		checkUser, err := s.repo.FindByUsername(ctx, req.Username)
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

		if err := s.repo.CreateUser(ctx, &newUser); err != nil {
			return err
		}

		newEmp := Employee{
			UserID:       newUser.ID,
			FullName:     req.FullName,
			NIK:          req.NIK,
			DepartmentID: req.DepartmentID,
			ShiftID:      req.ShiftID,
			BaseSalary:   req.BaseSalary,
		}

		if err := s.repo.CreateEmployee(ctx, &newEmp); err != nil {
			return err
		}

		err = s.leaveGenerator.GenerateInitialBalance(ctx, newEmp.ID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) UpdateEmployee(ctx context.Context, id uint, req *UpdateEmployeeRequest) error {
	emp, err := s.repo.FindEmployeeByID(ctx, id)
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
	if req.BaseSalary > 0 {
		emp.BaseSalary = req.BaseSalary
	}

	return s.repo.UpdateEmployee(ctx, emp)
}

func (s *service) DeleteEmployee(ctx context.Context, id uint) error {
	emp, err := s.repo.FindEmployeeByID(ctx, id)
	if err != nil {
		return errors.New("employee not found")
	}

	return s.repo.DeleteUser(ctx, emp.UserID)
}

func (s *service) buildEmployeeData(ctx context.Context, user *User, req *UpdateProfileRequest, file *multipart.FileHeader) (*User, error) {
	if req.FullName != "" {
		user.Employee.FullName = req.FullName
	}

	if req.PhoneNumber != "" {
		user.Employee.PhoneNumber = req.PhoneNumber
	}

	if req.BankName != "" {
		user.Employee.BankName = req.BankName
	}

	if req.BankAccountNumber != "" {
		user.Employee.BankAccountNumber = req.BankAccountNumber
	}

	if req.BankAccountHolder != "" {
		user.Employee.BankAccountHolder = req.BankAccountHolder
	}

	if req.NPWP != "" {
		user.Employee.NPWP = req.NPWP
	}

	if file != nil {
		fileName := fmt.Sprintf("users/%d/profile-%d.jpg", user.ID, time.Now().Unix())
		fileURL, err := s.storage.UploadFileMultipart(ctx, file, fileName)
		if err != nil {
			return nil, err
		}

		user.Employee.ProfilePictureUrl = fileURL
	}

	return user, nil
}
