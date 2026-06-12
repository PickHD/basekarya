package asset

import (
	"context"

	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/mock"
	"github.com/xuri/excelize/v2"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) CreateCategory(ctx context.Context, category *AssetCategory) error {
	return m.Called(ctx, category).Error(0)
}

func (m *mockRepo) FindCategoryByID(ctx context.Context, id uint) (*AssetCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AssetCategory), args.Error(1)
}

func (m *mockRepo) FindAllCategories(ctx context.Context, filter AssetCategoryFilter) ([]AssetCategory, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]AssetCategory), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) UpdateCategory(ctx context.Context, category *AssetCategory) error {
	return m.Called(ctx, category).Error(0)
}

func (m *mockRepo) DeleteCategory(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockRepo) CreateAsset(ctx context.Context, asset *Asset) error {
	return m.Called(ctx, asset).Error(0)
}

func (m *mockRepo) FindAssetByID(ctx context.Context, id uint) (*Asset, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Asset), args.Error(1)
}

func (m *mockRepo) FindAllAssets(ctx context.Context, filter AssetFilter) ([]Asset, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]Asset), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) UpdateAsset(ctx context.Context, asset *Asset) error {
	return m.Called(ctx, asset).Error(0)
}

func (m *mockRepo) DeleteAsset(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockRepo) CreateAssignment(ctx context.Context, assignment *AssetAssignment) error {
	return m.Called(ctx, assignment).Error(0)
}

func (m *mockRepo) FindAssignmentByID(ctx context.Context, id uint) (*AssetAssignment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AssetAssignment), args.Error(1)
}

func (m *mockRepo) FindActiveAssignmentByAssetID(ctx context.Context, assetID uint) (*AssetAssignment, error) {
	args := m.Called(ctx, assetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AssetAssignment), args.Error(1)
}

func (m *mockRepo) FindAllAssignments(ctx context.Context, filter AssetAssignmentFilter) ([]AssetAssignment, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]AssetAssignment), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) UpdateAssignment(ctx context.Context, assignment *AssetAssignment) error {
	return m.Called(ctx, assignment).Error(0)
}

type mockNotification struct{ mock.Mock }

func (m *mockNotification) SendNotification(ctx context.Context, userID uint, notifType string, title string, message string, relatedID uint) error {
	return m.Called(ctx, userID, notifType, title, message, relatedID).Error(0)
}

func (m *mockNotification) BlastNotification(ctx context.Context, userIDs []uint, notifType string, title string, message string, relatedID uint) error {
	return m.Called(ctx, userIDs, notifType, title, message, relatedID).Error(0)
}

type mockUserProvider struct{ mock.Mock }

func (m *mockUserProvider) FindApprovalUsers(ctx context.Context, permissionApprovalName string) ([]uint, error) {
	args := m.Called(ctx, permissionApprovalName)
	return args.Get(0).([]uint), args.Error(1)
}

type mockExcel struct{ mock.Mock }

func (m *mockExcel) GenerateSimpleExcel(sheetName string, headers []string, rows [][]interface{}) ([]byte, error) {
	args := m.Called(sheetName, headers, rows)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockExcel) NewFile() *excelize.File {
	args := m.Called()
	if args.Get(0) == nil {
		return excelize.NewFile()
	}
	return args.Get(0).(*excelize.File)
}

func (m *mockExcel) WriteToBuffer(file *excelize.File) ([]byte, error) {
	args := m.Called(file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

type mockService struct{ mock.Mock }

func (m *mockService) CreateCategory(ctx context.Context, req *CreateAssetCategoryRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetCategoryDetail(ctx context.Context, id uint) (*AssetCategoryResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AssetCategoryResponse), args.Error(1)
}

func (m *mockService) GetCategories(ctx context.Context, filter AssetCategoryFilter) ([]AssetCategoryResponse, *response.Meta, error) {
	args := m.Called(ctx, filter)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	return args.Get(0).([]AssetCategoryResponse), meta, args.Error(2)
}

func (m *mockService) UpdateCategory(ctx context.Context, req *UpdateAssetCategoryRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) DeleteCategory(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockService) CreateAsset(ctx context.Context, req *CreateAssetRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetAssetDetail(ctx context.Context, id uint) (*AssetDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AssetDetailResponse), args.Error(1)
}

func (m *mockService) GetAssets(ctx context.Context, filter AssetFilter) ([]AssetListResponse, *response.Meta, error) {
	args := m.Called(ctx, filter)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	return args.Get(0).([]AssetListResponse), meta, args.Error(2)
}

func (m *mockService) UpdateAsset(ctx context.Context, req *UpdateAssetRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) DeleteAsset(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockService) CreateAssignment(ctx context.Context, req *CreateAssetAssignmentRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetAssignmentDetail(ctx context.Context, id uint) (*AssetAssignmentDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AssetAssignmentDetailResponse), args.Error(1)
}

func (m *mockService) GetAssignments(ctx context.Context, filter AssetAssignmentFilter) ([]AssetAssignmentListResponse, *response.Meta, error) {
	args := m.Called(ctx, filter)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	return args.Get(0).([]AssetAssignmentListResponse), meta, args.Error(2)
}

func (m *mockService) ProcessAction(ctx context.Context, req *ActionRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) ProcessReturn(ctx context.Context, req *ReturnRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) Export(ctx context.Context, filter AssetFilter) ([]byte, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func approvalAssetKey() string { return string(constants.APPROVAL_ASSET) }
