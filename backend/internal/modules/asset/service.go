package asset

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"
	"basekarya-backend/pkg/utils"
	"context"
	"fmt"
)

type Service interface {
	CreateCategory(ctx context.Context, req *CreateAssetCategoryRequest) error
	GetCategoryDetail(ctx context.Context, id uint) (*AssetCategoryResponse, error)
	GetCategories(ctx context.Context, filter AssetCategoryFilter) ([]AssetCategoryResponse, *response.Meta, error)
	UpdateCategory(ctx context.Context, req *UpdateAssetCategoryRequest) error
	DeleteCategory(ctx context.Context, id uint) error

	CreateAsset(ctx context.Context, req *CreateAssetRequest) error
	GetAssetDetail(ctx context.Context, id uint) (*AssetDetailResponse, error)
	GetAssets(ctx context.Context, filter AssetFilter) ([]AssetListResponse, *response.Meta, error)
	UpdateAsset(ctx context.Context, req *UpdateAssetRequest) error
	DeleteAsset(ctx context.Context, id uint) error

	CreateAssignment(ctx context.Context, req *CreateAssetAssignmentRequest) error
	GetAssignmentDetail(ctx context.Context, id uint) (*AssetAssignmentDetailResponse, error)
	GetAssignments(ctx context.Context, filter AssetAssignmentFilter) ([]AssetAssignmentListResponse, *response.Meta, error)
	ProcessAction(ctx context.Context, req *ActionRequest) error
	ProcessReturn(ctx context.Context, req *ReturnRequest) error

	Export(ctx context.Context, filter AssetFilter) ([]byte, error)
}

type service struct {
	repo               Repository
	notification       NotificationProvider
	user               UserProvider
	transactionManager infrastructure.TransactionManager
	excel              infrastructure.ExcelProvider
}

func NewService(repo Repository, notification NotificationProvider, user UserProvider, transactionManager infrastructure.TransactionManager, excel infrastructure.ExcelProvider) Service {
	return &service{repo, notification, user, transactionManager, excel}
}

func (s *service) CreateCategory(ctx context.Context, req *CreateAssetCategoryRequest) error {
	category := &AssetCategory{
		CompanyID:   utils.GetCompanyIDFromCtx(ctx),
		Name:        req.Name,
		Description: req.Description,
	}
	return s.repo.CreateCategory(ctx, category)
}

func (s *service) GetCategoryDetail(ctx context.Context, id uint) (*AssetCategoryResponse, error) {
	category, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &AssetCategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (s *service) GetCategories(ctx context.Context, filter AssetCategoryFilter) ([]AssetCategoryResponse, *response.Meta, error) {
	categories, total, err := s.repo.FindAllCategories(ctx, filter)
	if err != nil {
		return []AssetCategoryResponse{}, nil, nil
	}

	if len(categories) == 0 {
		return []AssetCategoryResponse{}, nil, nil
	}

	var list []AssetCategoryResponse
	for _, c := range categories {
		list = append(list, AssetCategoryResponse{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		})
	}

	meta := response.NewMetaOffset(filter.Page, filter.Limit, total)
	return list, meta, nil
}

func (s *service) UpdateCategory(ctx context.Context, req *UpdateAssetCategoryRequest) error {
	category, err := s.repo.FindCategoryByID(ctx, req.ID)
	if err != nil {
		return err
	}
	category.Name = req.Name
	category.Description = req.Description
	return s.repo.UpdateCategory(ctx, category)
}

func (s *service) DeleteCategory(ctx context.Context, id uint) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *service) CreateAsset(ctx context.Context, req *CreateAssetRequest) error {
	condition := constants.AssetConditionGood
	if req.Condition != "" {
		condition = req.Condition
	}

	asset := &Asset{
		CompanyID:       utils.GetCompanyIDFromCtx(ctx),
		AssetCategoryID: req.AssetCategoryID,
		Name:            req.Name,
		Description:     req.Description,
		SerialNumber:    req.SerialNumber,
		Status:          constants.AssetStatusAvailable,
		Condition:       condition,
	}
	return s.repo.CreateAsset(ctx, asset)
}

func (s *service) GetAssetDetail(ctx context.Context, id uint) (*AssetDetailResponse, error) {
	asset, err := s.repo.FindAssetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	categoryName := ""
	if asset.AssetCategory.ID != 0 {
		categoryName = asset.AssetCategory.Name
	}

	currentEmployee := ""
	activeAssignment, err := s.repo.FindActiveAssignmentByAssetID(ctx, id)
	if err == nil && activeAssignment != nil && activeAssignment.Employee.ID != 0 {
		currentEmployee = activeAssignment.Employee.FullName
	}

	return &AssetDetailResponse{
		ID:              asset.ID,
		Name:            asset.Name,
		Description:     asset.Description,
		SerialNumber:    asset.SerialNumber,
		AssetCategoryID: asset.AssetCategoryID,
		CategoryName:    categoryName,
		Status:          asset.Status,
		Condition:       asset.Condition,
		CurrentEmployee: currentEmployee,
		CreatedAt:       asset.CreatedAt,
		UpdatedAt:       asset.UpdatedAt,
	}, nil
}

func (s *service) GetAssets(ctx context.Context, filter AssetFilter) ([]AssetListResponse, *response.Meta, error) {
	assets, total, err := s.repo.FindAllAssets(ctx, filter)
	if err != nil {
		return []AssetListResponse{}, nil, nil
	}

	if len(assets) == 0 {
		return []AssetListResponse{}, nil, nil
	}

	var list []AssetListResponse
	for _, a := range assets {
		categoryName := ""
		if a.AssetCategory.ID != 0 {
			categoryName = a.AssetCategory.Name
		}
		list = append(list, AssetListResponse{
			ID:              a.ID,
			Name:            a.Name,
			Description:     a.Description,
			SerialNumber:    a.SerialNumber,
			AssetCategoryID: a.AssetCategoryID,
			CategoryName:    categoryName,
			Status:          a.Status,
			Condition:       a.Condition,
			CreatedAt:       a.CreatedAt,
			UpdatedAt:       a.UpdatedAt,
		})
	}

	meta := response.NewMetaOffset(filter.Page, filter.Limit, total)
	return list, meta, nil
}

func (s *service) UpdateAsset(ctx context.Context, req *UpdateAssetRequest) error {
	asset, err := s.repo.FindAssetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if req.AssetCategoryID > 0 {
		asset.AssetCategoryID = req.AssetCategoryID
	}
	if req.Name != "" {
		asset.Name = req.Name
	}
	if req.Description != "" {
		asset.Description = req.Description
	}
	if req.SerialNumber != "" {
		asset.SerialNumber = req.SerialNumber
	}
	if req.Status != "" {
		asset.Status = req.Status
	}
	if req.Condition != "" {
		asset.Condition = req.Condition
	}

	return s.repo.UpdateAsset(ctx, asset)
}

func (s *service) DeleteAsset(ctx context.Context, id uint) error {
	return s.repo.DeleteAsset(ctx, id)
}

func (s *service) CreateAssignment(ctx context.Context, req *CreateAssetAssignmentRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		if req.UserID == 0 || req.EmployeeID == 0 {
			return fmt.Errorf("user not found")
		}

		asset, err := s.repo.FindAssetByID(ctx, req.AssetID)
		if err != nil {
			return fmt.Errorf("asset not found")
		}

		if asset.Status != constants.AssetStatusAvailable {
			return fmt.Errorf("asset is not available for assignment")
		}

		assignment := &AssetAssignment{
			CompanyID:  utils.GetCompanyIDFromCtx(ctx),
			AssetID:    req.AssetID,
			EmployeeID: req.EmployeeID,
			UserID:     req.UserID,
			Purpose:    req.Purpose,
			Status:     constants.AssetAssignmentStatusPending,
		}
		if req.ExpectedReturnDate != "" {
			assignment.ExpectedReturnDate = &req.ExpectedReturnDate
		}

		err = s.repo.CreateAssignment(ctx, assignment)
		if err != nil {
			return err
		}

		approvalUserIDs, err := s.user.FindApprovalUsers(ctx, string(constants.APPROVAL_ASSET))
		if err != nil {
			return err
		}

		go func() {
			_ = s.notification.BlastNotification(
				utils.DetachContext(ctx),
				approvalUserIDs,
				string(constants.NotificationTypeAssetApprovalReq),
				"Permintaan Aset Baru",
				fmt.Sprintf("Karyawan mengajukan permintaan aset %s", asset.Name),
				assignment.ID,
			)
		}()

		return nil
	})
}

func (s *service) GetAssignmentDetail(ctx context.Context, id uint) (*AssetAssignmentDetailResponse, error) {
	assignment, err := s.repo.FindAssignmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if assignment.User.ID == 0 && assignment.Employee.ID == 0 {
		return nil, fmt.Errorf("data user not found")
	}

	rejectionReason := ""
	if assignment.RejectionReason.Valid {
		rejectionReason = assignment.RejectionReason.String
	}

	assetName := ""
	if assignment.Asset.ID != 0 {
		assetName = assignment.Asset.Name
	}

	return &AssetAssignmentDetailResponse{
		ID:                 assignment.ID,
		AssetID:            assignment.AssetID,
		AssetName:          assetName,
		EmployeeID:         assignment.EmployeeID,
		EmployeeName:       assignment.Employee.FullName,
		EmployeeNIK:        assignment.Employee.NIK,
		Purpose:            assignment.Purpose,
		ExpectedReturnDate: assignment.ExpectedReturnDate,
		ActualReturnDate:   assignment.ActualReturnDate,
		Notes:              assignment.Notes,
		Status:             assignment.Status,
		RejectionReason:    rejectionReason,
		CreatedAt:          assignment.CreatedAt,
	}, nil
}

func (s *service) GetAssignments(ctx context.Context, filter AssetAssignmentFilter) ([]AssetAssignmentListResponse, *response.Meta, error) {
	assignments, total, err := s.repo.FindAllAssignments(ctx, filter)
	if err != nil {
		return []AssetAssignmentListResponse{}, nil, nil
	}

	if len(assignments) == 0 {
		return []AssetAssignmentListResponse{}, nil, nil
	}

	var list []AssetAssignmentListResponse
	for _, a := range assignments {
		assetName := ""
		if a.Asset.ID != 0 {
			assetName = a.Asset.Name
		}

		list = append(list, AssetAssignmentListResponse{
			ID:                 a.ID,
			AssetID:            a.AssetID,
			AssetName:          assetName,
			EmployeeID:         a.EmployeeID,
			EmployeeName:       a.Employee.FullName,
			EmployeeNIK:        a.Employee.NIK,
			Purpose:            a.Purpose,
			ExpectedReturnDate: a.ExpectedReturnDate,
			ActualReturnDate:   a.ActualReturnDate,
			Status:             a.Status,
			CreatedAt:          a.CreatedAt,
		})
	}

	meta := response.NewMetaOffset(filter.Page, filter.Limit, total)
	return list, meta, nil
}

func (s *service) ProcessAction(ctx context.Context, req *ActionRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		data, err := s.repo.FindAssignmentByID(ctx, req.ID)
		if err != nil {
			return err
		}

		if data.Status != constants.AssetAssignmentStatusPending {
			return fmt.Errorf("cannot process assignment with status %s", data.Status)
		}

		var (
			notificationType    constants.NotificationType
			notificationTitle   string
			notificationMessage string
		)
		switch constants.AssetAssignmentAction(req.Action) {
		case constants.AssetAssignmentActionApprove:
			data.Status = constants.AssetAssignmentStatusActive
			data.ApprovedBy = &req.SuperAdminID

			err = s.repo.UpdateAssignment(ctx, data)
			if err != nil {
				return err
			}

			asset, err := s.repo.FindAssetByID(ctx, data.AssetID)
			if err != nil {
				return err
			}
			asset.Status = constants.AssetStatusAssigned
			err = s.repo.UpdateAsset(ctx, asset)
			if err != nil {
				return err
			}

			notificationType = constants.NotificationTypeApproved
			notificationTitle = "Permintaan Disetujui"
			notificationMessage = "Permintaan aset Anda telah disetujui oleh Admin."
		case constants.AssetAssignmentActionReject:
			if req.RejectionReason == "" {
				return fmt.Errorf("rejection reason is required")
			}
			data.Status = constants.AssetAssignmentStatusRejected
			data.RejectionReason.String = req.RejectionReason
			data.RejectionReason.Valid = true

			err = s.repo.UpdateAssignment(ctx, data)
			if err != nil {
				return err
			}

			notificationType = constants.NotificationTypeRejected
			notificationTitle = "Permintaan Ditolak"
			notificationMessage = "Permintaan aset Anda telah ditolak oleh Admin."
		default:
			return fmt.Errorf("invalid action: %s", req.Action)
		}

		go func() {
			_ = s.notification.SendNotification(
				utils.DetachContext(ctx),
				data.UserID,
				string(notificationType),
				notificationTitle,
				notificationMessage,
				data.ID,
			)
		}()

		return nil
	})
}

func (s *service) ProcessReturn(ctx context.Context, req *ReturnRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		data, err := s.repo.FindAssignmentByID(ctx, req.ID)
		if err != nil {
			return err
		}

		if data.Status != constants.AssetAssignmentStatusActive {
			return fmt.Errorf("cannot return assignment with status %s", data.Status)
		}

		data.Status = constants.AssetAssignmentStatusReturned
		err = s.repo.UpdateAssignment(ctx, data)
		if err != nil {
			return err
		}

		asset, err := s.repo.FindAssetByID(ctx, data.AssetID)
		if err != nil {
			return err
		}
		asset.Status = constants.AssetStatusAvailable
		return s.repo.UpdateAsset(ctx, asset)
	})
}

func (s *service) Export(ctx context.Context, filter AssetFilter) ([]byte, error) {
	filter.Page = 1
	filter.Limit = 999999

	assets, _, err := s.repo.FindAllAssets(ctx, filter)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"ID", "Nama", "Deskripsi", "Serial Number", "Kategori", "Status", "Kondisi", "Dibuat Pada",
	}

	var rows [][]interface{}
	for _, a := range assets {
		categoryName := "-"
		if a.AssetCategory.ID != 0 {
			categoryName = a.AssetCategory.Name
		}

		row := []interface{}{
			a.ID,
			a.Name,
			a.Description,
			a.SerialNumber,
			categoryName,
			a.Status,
			a.Condition,
			a.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		rows = append(rows, row)
	}

	return s.excel.GenerateSimpleExcel("Assets", headers, rows)
}


