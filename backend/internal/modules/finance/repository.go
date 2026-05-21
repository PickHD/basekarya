package finance

import (
	"basekarya-backend/pkg/utils"
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	CreateTransaction(ctx context.Context, tx *FinanceTransaction) error
	FindTransactionByID(ctx context.Context, id uint) (*FinanceTransaction, error)
	FindAllTransactions(ctx context.Context, filter TransactionFilter) ([]FinanceTransaction, int64, error)
	UpdateTransaction(ctx context.Context, tx *FinanceTransaction) error

	CreateCategory(ctx context.Context, cat *FinanceCategory) error
	FindCategoryByID(ctx context.Context, id uint) (*FinanceCategory, error)
	FindAllCategories(ctx context.Context, catType string) ([]FinanceCategory, error)
	UpdateCategory(ctx context.Context, cat *FinanceCategory) error
	DeleteCategory(ctx context.Context, id uint) error

	GetDashboardSummary(ctx context.Context, startDate, endDate string) (*DashboardResponse, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateTransaction(ctx context.Context, tx *FinanceTransaction) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Create(tx).Error
}

func (r *repository) FindTransactionByID(ctx context.Context, id uint) (*FinanceTransaction, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var tx FinanceTransaction

	err := db.
		Preload("Creator").
		Preload("Creator.Employee").
		Preload("FinanceCategory").
		Preload("Approver").
		Preload("Approver.Employee").
		First(&tx, id).Error
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func (r *repository) FindAllTransactions(ctx context.Context, filter TransactionFilter) ([]FinanceTransaction, int64, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var transactions []FinanceTransaction
	var total int64

	query := db.Model(&FinanceTransaction{}).
		Joins("JOIN users ON users.id = finance_transactions.created_by").
		Joins("JOIN finance_categories ON finance_categories.id = finance_transactions.finance_category_id").
		Preload("Creator").
		Preload("Creator.Employee").
		Preload("FinanceCategory")

	if filter.CreatedBy > 0 {
		query = query.Where("finance_transactions.created_by = ?", filter.CreatedBy)
	}

	if filter.Type != "" {
		query = query.Where("finance_transactions.type = ?", filter.Type)
	}

	if filter.Status != "" {
		query = query.Where("finance_transactions.status = ?", filter.Status)
	}

	if filter.StartDate != "" {
		query = query.Where("finance_transactions.transaction_date >= ?", filter.StartDate)
	}

	if filter.EndDate != "" {
		query = query.Where("finance_transactions.transaction_date <= ?", filter.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Limit(filter.Limit).
		Offset(offset).
		Order("finance_transactions.created_at DESC").
		Find(&transactions).Error

	return transactions, total, err
}

func (r *repository) UpdateTransaction(ctx context.Context, tx *FinanceTransaction) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Save(tx).Error
}

func (r *repository) CreateCategory(ctx context.Context, cat *FinanceCategory) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Create(cat).Error
}

func (r *repository) FindCategoryByID(ctx context.Context, id uint) (*FinanceCategory, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var cat FinanceCategory
	err := db.First(&cat, id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *repository) FindAllCategories(ctx context.Context, catType string) ([]FinanceCategory, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var categories []FinanceCategory

	query := db.Order("name ASC")
	if catType != "" {
		query = query.Where("type = ?", catType)
	}

	err := query.Find(&categories).Error
	return categories, err
}

func (r *repository) UpdateCategory(ctx context.Context, cat *FinanceCategory) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Save(cat).Error
}

func (r *repository) DeleteCategory(ctx context.Context, id uint) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Delete(&FinanceCategory{}, id).Error
}

func (r *repository) GetDashboardSummary(ctx context.Context, startDate, endDate string) (*DashboardResponse, error) {
	db := utils.GetDBFromContext(ctx, r.db)

	resp := &DashboardResponse{}

	statusFilter := "APPROVED"

	incomeQuery := db.Model(&FinanceTransaction{}).
		Where("type = ? AND status = ?", "INCOME", statusFilter)
	expenseQuery := db.Model(&FinanceTransaction{}).
		Where("type = ? AND status = ?", "EXPENSE", statusFilter)

	if startDate != "" {
		incomeQuery = incomeQuery.Where("transaction_date >= ?", startDate)
		expenseQuery = expenseQuery.Where("transaction_date >= ?", startDate)
	}
	if endDate != "" {
		incomeQuery = incomeQuery.Where("transaction_date <= ?", endDate)
		expenseQuery = expenseQuery.Where("transaction_date <= ?", endDate)
	}

	var totalIncome, totalExpense float64
	incomeQuery.Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)
	expenseQuery.Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense)

	resp.TotalIncome = totalIncome
	resp.TotalExpense = totalExpense
	resp.NetBalance = totalIncome - totalExpense

	baseCountQuery := db.Model(&FinanceTransaction{}).Where("status = ?", statusFilter)
	if startDate != "" {
		baseCountQuery = baseCountQuery.Where("transaction_date >= ?", startDate)
	}
	if endDate != "" {
		baseCountQuery = baseCountQuery.Where("transaction_date <= ?", endDate)
	}
	var count int64
	baseCountQuery.Count(&count)
	resp.TransactionCount = count

	type monthResult struct {
		Month  string
		Total  float64
		TxType string
	}

	var monthResults []monthResult
	monthQuery := db.Model(&FinanceTransaction{}).
		Select("DATE_FORMAT(transaction_date, '%Y-%m') as month, SUM(amount) as total, type as tx_type").
		Where("status = ?", statusFilter)

	if startDate != "" {
		monthQuery = monthQuery.Where("transaction_date >= ?", startDate)
	}
	if endDate != "" {
		monthQuery = monthQuery.Where("transaction_date <= ?", endDate)
	}

	monthQuery.Group("month, tx_type").Order("month ASC").Find(&monthResults)

	monthMap := make(map[string]*MonthlySummaryItem)
	for _, mr := range monthResults {
		if _, ok := monthMap[mr.Month]; !ok {
			monthMap[mr.Month] = &MonthlySummaryItem{Month: mr.Month}
		}
		if mr.TxType == "INCOME" {
			monthMap[mr.Month].Income = mr.Total
		} else {
			monthMap[mr.Month].Expense = mr.Total
		}
	}

	var monthlySummary []MonthlySummaryItem
	for _, item := range monthMap {
		monthlySummary = append(monthlySummary, *item)
	}
	resp.MonthlySummary = monthlySummary

	type catResult struct {
		CategoryName string
		TxType       string
		Total        float64
	}

	var catResults []catResult
	catQuery := db.Model(&FinanceTransaction{}).
		Select("finance_categories.name as category_name, finance_transactions.type as tx_type, SUM(finance_transactions.amount) as total").
		Joins("JOIN finance_categories ON finance_categories.id = finance_transactions.finance_category_id").
		Where("finance_transactions.status = ?", statusFilter)

	if startDate != "" {
		catQuery = catQuery.Where("finance_transactions.transaction_date >= ?", startDate)
	}
	if endDate != "" {
		catQuery = catQuery.Where("finance_transactions.transaction_date <= ?", endDate)
	}

	catQuery.Group("finance_categories.name, finance_transactions.type").Find(&catResults)

	var categoryBreakdown []CategoryBreakdownItem
	for _, cr := range catResults {
		categoryBreakdown = append(categoryBreakdown, CategoryBreakdownItem{
			CategoryName: cr.CategoryName,
			Type:         cr.TxType,
			Total:        cr.Total,
		})
	}
	resp.CategoryBreakdown = categoryBreakdown

	var recentTx []FinanceTransaction
	db.Preload("Creator").
		Preload("Creator.Employee").
		Preload("FinanceCategory").
		Where("status = ?", statusFilter).
		Order("created_at DESC").
		Limit(5).
		Find(&recentTx)

	var recentResponse []TransactionListResponse
	for _, tx := range recentTx {
		creatorName := ""
		if tx.Creator.Employee != nil {
			creatorName = tx.Creator.Employee.FullName
		}

		recentResponse = append(recentResponse, TransactionListResponse{
			ID:              tx.ID,
			CreatorName:     creatorName,
			CategoryName:    tx.FinanceCategory.Name,
			Type:            tx.Type,
			Amount:          tx.Amount,
			TransactionDate: tx.TransactionDate,
			ReferenceNumber: tx.ReferenceNumber.String,
			Status:          tx.Status,
			CreatedAt:       tx.CreatedAt,
		})
	}
	resp.RecentTransactions = recentResponse

	return resp, nil
}
