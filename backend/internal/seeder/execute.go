package seeder

import (
	"errors"
	"strings"

	"basekarya-backend/internal/config"
	"basekarya-backend/internal/modules/department"
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/subscription"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/logger"

	"gorm.io/gorm"
)

func Execute(db *gorm.DB, cfg *config.Config, hasher Hasher) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := seedMasterData(tx); err != nil {
			return err
		}

		roleSuperadmin, err := seedRolesAndAdmin(tx, cfg, hasher)
		if err != nil {
			return err
		}

		if err := seedPermissions(tx, roleSuperadmin); err != nil {
			return err
		}

		if err := seedSubscriptionPlans(tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Errorf("Seeding failed: %v", err)
		return err
	}

	logger.Info("Database seeding completed successfully!")
	return nil
}

func seedMasterData(tx *gorm.DB) error {
	generalDept := department.Department{Name: "Umum", CompanyID: 1}
	if err := tx.Where(department.Department{Name: "Umum"}).FirstOrCreate(&generalDept).Error; err != nil {
		return err
	}

	regularShift := master.Shift{Name: "Regular", StartTime: "09:00:00", EndTime: "18:00:00", CompanyID: 1}
	var existingShift master.Shift
	err := tx.Where("name = ?", regularShift.Name).First(&existingShift).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.Create(&regularShift).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		if err := tx.Model(&existingShift).Updates(map[string]interface{}{
			"start_time": regularShift.StartTime,
			"end_time":   regularShift.EndTime,
		}).Error; err != nil {
			return err
		}
	}

	leaveTypes := []master.LeaveType{
		{Name: "Annual", DefaultQuota: 12, IsDeducted: true, CompanyID: 1},
		{Name: "Sick", DefaultQuota: 15, IsDeducted: false, CompanyID: 1},
		{Name: "Unpaid", DefaultQuota: 0, IsDeducted: false, CompanyID: 1},
	}

	for _, lt := range leaveTypes {
		var existing master.LeaveType
		err := tx.Where("name = ?", lt.Name).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&lt).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			if err := tx.Model(&existing).Updates(map[string]interface{}{
				"default_quota": lt.DefaultQuota,
				"is_deducted":   lt.IsDeducted,
			}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func seedRolesAndAdmin(tx *gorm.DB, cfg *config.Config, hasher Hasher) (*rbac.Role, error) {
	platformRole := rbac.Role{Name: string(constants.UserRolePlatformAdmin), CompanyID: 0}
	if err := tx.Where(rbac.Role{Name: string(constants.UserRolePlatformAdmin)}).FirstOrCreate(&platformRole).Error; err != nil {
		return nil, err
	}

	hashPass, err := hasher.HashPassword(cfg.CredentialConfig.SuperadminPassword)
	if err != nil {
		return nil, err
	}

	var existingAdmin user.User
	err = tx.Where("username = ?", cfg.CredentialConfig.SuperadminUsername).First(&existingAdmin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newAdmin := user.User{
			Username:           cfg.CredentialConfig.SuperadminUsername,
			PasswordHash:       hashPass,
			RoleID:             platformRole.ID,
			CompanyID:          0,
			IsPlatformAdmin:    true,
			MustChangePassword: false,
			IsActive:           true,
		}
		if err := tx.Create(&newAdmin).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		if err := tx.Model(&existingAdmin).Updates(map[string]interface{}{
			"password_hash":        hashPass,
			"role_id":              platformRole.ID,
			"company_id":           0,
			"is_platform_admin":    true,
			"must_change_password": false,
			"is_active":            true,
		}).Error; err != nil {
			return nil, err
		}
	}

	return &platformRole, nil
}

func seedPermissions(tx *gorm.DB, roleSuperadmin *rbac.Role) error {
	permissionsByGroup := []struct {
		GroupName   string
		Permissions []string
	}{
		{"Permission", []string{constants.VIEW_PERMISSION}},
		{"Role", []string{constants.CREATE_ROLE, constants.VIEW_ROLE, constants.ASSIGN_ROLE}},
		{"Master", []string{constants.VIEW_MASTER, constants.MANAGE_MASTER}},
		{"Employee", []string{constants.VIEW_EMPLOYEE, constants.CREATE_EMPLOYEE, constants.UPDATE_EMPLOYEE, constants.DELETE_EMPLOYEE, constants.EXPORT_EMPLOYEE}},
		{"Attendance", []string{constants.VIEW_ATTENDANCE, constants.VIEW_SELF_ATTENDANCE, constants.CREATE_ATTENDANCE, constants.EXPORT_ATTENDANCE}},
		{"Payroll", []string{constants.VIEW_PAYROLL, constants.GENERATE_PAYROLL, constants.DOWNLOAD_PAYSLIP, constants.MARK_AS_PAID, constants.SEND_PAYSLIP}},
		{"Leave", []string{constants.VIEW_LEAVE, constants.VIEW_SELF_LEAVE, constants.CREATE_LEAVE, constants.APPROVAL_LEAVE, constants.EXPORT_LEAVE}},
		{"Loan", []string{constants.VIEW_LOAN, constants.VIEW_SELF_LOAN, constants.CREATE_LOAN, constants.APPROVAL_LOAN, constants.EXPORT_LOAN}},
		{"Overtime", []string{constants.VIEW_OVERTIME, constants.VIEW_SELF_OVERTIME, constants.CREATE_OVERTIME, constants.APPROVAL_OVERTIME, constants.EXPORT_OVERTIME}},
		{"Reimbursement", []string{constants.VIEW_REIMBURSEMENT, constants.VIEW_SELF_REIMBURSEMENT, constants.CREATE_REIMBURSEMENT, constants.APPROVAL_REIMBURSEMENT, constants.EXPORT_REIMBURSEMENT}},
		{"Company", []string{constants.VIEW_COMPANY, constants.UPDATE_COMPANY}},
		{"Announcement", []string{constants.CREATE_ANNOUNCEMENT}},
		{"Contract", []string{constants.VIEW_CONTRACT, constants.CREATE_CONTRACT, constants.UPDATE_CONTRACT, constants.EXPORT_CONTRACT}},
		{"Recruitment", []string{constants.VIEW_REQUISITION, constants.CREATE_REQUISITION, constants.APPROVAL_REQUISITION, constants.VIEW_APPLICANT, constants.CREATE_APPLICANT, constants.UPDATE_APPLICANT}},
		{"Onboarding", []string{constants.VIEW_ONBOARDING, constants.MANAGE_ONBOARDING_TEMPLATE, constants.UPDATE_ONBOARDING_TASK}},
		{"Finance", []string{constants.VIEW_FINANCE, constants.CREATE_FINANCE, constants.APPROVAL_FINANCE, constants.EXPORT_FINANCE, constants.MANAGE_FINANCE_CATEGORY, constants.VIEW_FINANCE_DASHBOARD}},
		{"Asset", []string{constants.MANAGE_ASSET, constants.VIEW_ASSET, constants.VIEW_SELF_ASSET, constants.CREATE_ASSET, constants.APPROVAL_ASSET, constants.EXPORT_ASSET}},
	}

	var permissionIDs []uint

	for _, pg := range permissionsByGroup {
		group := rbac.PermissionGroup{Name: pg.GroupName}
		if err := tx.Where(rbac.PermissionGroup{Name: group.Name}).FirstOrCreate(&group).Error; err != nil {
			return err
		}

		for _, permName := range pg.Permissions {
			var perm rbac.Permission
			displayName := formatDisplayName(permName)

			if err := tx.Where(rbac.Permission{Name: permName}).
				Assign(rbac.Permission{
					PermissionGroupID: group.ID,
					DisplayName:       displayName,
				}).
				FirstOrCreate(&perm).Error; err != nil {
				return err
			}
			permissionIDs = append(permissionIDs, perm.ID)
		}
	}

	// Assign all permissions to platform admin role (company_id = 0)
	for _, pid := range permissionIDs {
		var rp rbac.RolePermission
		if err := tx.Where(rbac.RolePermission{RoleID: roleSuperadmin.ID, PermissionID: pid}).Attrs(rbac.RolePermission{CompanyID: 0}).FirstOrCreate(&rp).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedSubscriptionPlans(tx *gorm.DB) error {
	plans := []subscription.SubscriptionPlan{
		{Name: "Free", Slug: "free", MaxEmployees: 5, PriceMonthly: 0, Features: `{"modules":["attendance","leave"]}`, IsActive: true},
		{Name: "Basic", Slug: "basic", MaxEmployees: 50, PriceMonthly: 99000, Features: `{"modules":["attendance","leave","overtime","loan","reimbursement","payroll","contract","finance"]}`, IsActive: true},
		{Name: "Pro", Slug: "pro", MaxEmployees: 0, PriceMonthly: 249000, Features: `{"modules":["attendance","leave","overtime","loan","reimbursement","payroll","contract","finance","recruitment","onboarding","asset"]}`, IsActive: true},
	}

	for i := range plans {
		p := &plans[i]
		var existing subscription.SubscriptionPlan
		err := tx.Where("slug = ?", p.Slug).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(p).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			if err := tx.Model(&existing).Updates(map[string]interface{}{
				"name":           p.Name,
				"max_employees":  p.MaxEmployees,
				"price_monthly":  p.PriceMonthly,
				"features":       p.Features,
				"is_active":      p.IsActive,
			}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func formatDisplayName(name string) string {
	words := strings.Split(strings.ToLower(name), "_")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

