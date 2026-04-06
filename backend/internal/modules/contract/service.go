package contract

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"
	"basekarya-backend/pkg/utils"
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"
)

type Service interface {
	Upsert(ctx context.Context, req *UpsertContractRequest) error
	GetList(ctx context.Context, filter *ContractFilter) ([]ContractListResponse, *response.Meta, error)
	GetDetail(ctx context.Context, id uint) (*ContractDetailResponse, error)
	GetByEmployeeID(ctx context.Context, employeeID uint) (*ContractDetailResponse, error)
	Delete(ctx context.Context, id uint) error
	Export(ctx context.Context, filter *ContractFilter) ([]byte, error)
	CheckExpiringContracts(ctx context.Context) error
}

type service struct {
	repo         Repository
	storage      StorageProvider
	notification NotificationProvider
	user         UserProvider
	excel        infrastructure.ExcelProvider
}

func NewService(repo Repository, storage StorageProvider, notification NotificationProvider, user UserProvider, excel infrastructure.ExcelProvider) Service {
	return &service{repo, storage, notification, user, excel}
}

func (s *service) Upsert(ctx context.Context, req *UpsertContractRequest) error {
	start, err := time.Parse(constants.DefaultTimeFormat, req.StartDate)
	if err != nil {
		return errors.New("invalid start date format")
	}

	var end *time.Time
	if req.ContractType == constants.ContractTypePKWT {
		if req.EndDate == "" {
			return errors.New("end date is required for PKWT")
		}
		e, err := time.Parse(constants.DefaultTimeFormat, req.EndDate)
		if err != nil {
			return errors.New("invalid end date format")
		}
		if e.Before(start) {
			return errors.New("end date must be after start date")
		}
		end = &e
	}

	// Fetch existing to see if endDate changed
	existing, _ := s.repo.FindByEmployeeID(ctx, req.EmployeeID)

	attachmentUrl := ""
	if existing != nil {
		attachmentUrl = existing.AttachmentURL
	}

	if req.AttachmentBase64 != "" {
		imgBytes, err := utils.DecodeBase64Image(req.AttachmentBase64)
		if err != nil {
			return errors.New("invalid attachment")
		}
		imageReader := bytes.NewReader(imgBytes)
		now := time.Now()
		fileName := fmt.Sprintf("contracts/%d/%s.jpg", req.EmployeeID, now.Format("20060102_150405"))

		attachmentUrl, err = s.storage.UploadFileByte(ctx, fileName, imageReader, int64(len(imgBytes)), "image/jpeg")
		if err != nil {
			return err
		}
	}

	contract := &Contract{
		EmployeeID:     req.EmployeeID,
		ContractType:   req.ContractType,
		ContractNumber: req.ContractNumber,
		StartDate:      start,
		EndDate:        end,
		Notes:          req.Notes,
		AttachmentURL:  attachmentUrl,
	}

	// Reset alerted_at if end date is extended/changed
	if existing != nil && existing.AlertedAt != nil {
		if existing.EndDate != nil && end != nil && !existing.EndDate.Equal(*end) {
			contract.AlertedAt = nil
		} else {
			contract.AlertedAt = existing.AlertedAt
		}
	}

	return s.repo.Upsert(ctx, contract)
}

func (s *service) GetList(ctx context.Context, filter *ContractFilter) ([]ContractListResponse, *response.Meta, error) {
	contracts, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	var list []ContractListResponse
	for _, c := range contracts {
		employeeName := "-"
		employeeNIK := "-"
		if c.Employee != nil {
			employeeName = c.Employee.FullName
			employeeNIK = c.Employee.NIK
		}
		list = append(list, ContractListResponse{
			ID:             c.ID,
			EmployeeID:     c.EmployeeID,
			EmployeeName:   employeeName,
			EmployeeNIK:    employeeNIK,
			ContractType:   c.ContractType,
			ContractNumber: c.ContractNumber,
			StartDate:      c.StartDate,
			EndDate:        c.EndDate,
			CreatedAt:      c.CreatedAt,
		})
	}

	meta := response.NewMetaOffset(filter.Page, filter.Limit, total)
	return list, meta, nil
}

func (s *service) GetDetail(ctx context.Context, id uint) (*ContractDetailResponse, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapToDetailResponse(c), nil
}

func (s *service) GetByEmployeeID(ctx context.Context, employeeID uint) (*ContractDetailResponse, error) {
	c, err := s.repo.FindByEmployeeID(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	return s.mapToDetailResponse(c), nil
}

func (s *service) mapToDetailResponse(c *Contract) *ContractDetailResponse {
	employeeName := "-"
	employeeNIK := "-"
	if c.Employee != nil {
		employeeName = c.Employee.FullName
		employeeNIK = c.Employee.NIK
	}
	return &ContractDetailResponse{
		ID:             c.ID,
		EmployeeID:     c.EmployeeID,
		EmployeeName:   employeeName,
		EmployeeNIK:    employeeNIK,
		ContractType:   c.ContractType,
		ContractNumber: c.ContractNumber,
		StartDate:      c.StartDate,
		EndDate:        c.EndDate,
		Notes:          c.Notes,
		AttachmentURL:  c.AttachmentURL,
		CreatedAt:      c.CreatedAt,
	}
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *service) Export(ctx context.Context, filter *ContractFilter) ([]byte, error) {
	filter.Page = 1
	filter.Limit = 999999
	contracts, _, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"Nama Karyawan", "NIK", "Tipe Kontrak", "Nomor Kontrak", "Mulai", "Selesai",
	}

	var rows [][]interface{}
	for _, req := range contracts {
		empName := "-"
		empNik := "-"
		if req.Employee != nil {
			empName = req.Employee.FullName
			empNik = req.Employee.NIK
		}
		
		endDate := "-"
		if req.EndDate != nil {
			endDate = req.EndDate.Format(constants.DefaultTimeFormat)
		}

		row := []interface{}{
			empName,
			empNik,
			req.ContractType,
			req.ContractNumber,
			req.StartDate.Format(constants.DefaultTimeFormat),
			endDate,
		}
		rows = append(rows, row)
	}

	return s.excel.GenerateSimpleExcel("Contracts", headers, rows)
}

func (s *service) CheckExpiringContracts(ctx context.Context) error {
	contracts, err := s.repo.FindExpiringContracts(ctx, 30) // Expiring in <= 30 days
	if err != nil {
		return err
	}

	if len(contracts) == 0 {
		return nil
	}

	approvalUserIDs, err := s.user.FindApprovalUsers(ctx, "VIEW_CONTRACT")
	if err != nil || len(approvalUserIDs) == 0 {
		logger.Warn("CheckExpiringContracts: No users found with VIEW_CONTRACT permission")
		return nil
	}

	var notifiedIDs []uint
	for _, c := range contracts {
		if c.Employee == nil {
			contractWithEmp, err := s.repo.FindByID(ctx, c.ID)
			if err == nil {
				c = *contractWithEmp
			}
		}

		empName := ""
		if c.Employee != nil {
			empName = c.Employee.FullName
		}

		message := fmt.Sprintf("Kontrak PKWT atas nama %s (%s) akan segera berakhir pada %s.", empName, c.ContractNumber, c.EndDate.Format(constants.DefaultTimeFormat))

		_ = s.notification.BlastNotification(
			approvalUserIDs,
			"CONTRACT_EXPIRING",
			"Kontrak Segera Berakhir",
			message,
			c.ID,
		)

		notifiedIDs = append(notifiedIDs, c.ID)
	}

	if len(notifiedIDs) > 0 {
		_ = s.repo.MarkAlerted(ctx, notifiedIDs)
	}

	return nil
}
