package seeder

import (
	"strings"

	"basekarya-backend/internal/config"
	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/rbac"
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

		if err := seedCompanyData(tx); err != nil {
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
	generalDept := master.Department{Name: "Umum"}
	if err := tx.Where(master.Department{Name: "Umum"}).FirstOrCreate(&generalDept).Error; err != nil {
		return err
	}

	regularShift := master.Shift{Name: "Regular", StartTime: "09:00:00", EndTime: "18:00:00"}
	if err := tx.FirstOrCreate(&regularShift, master.Shift{Name: "Regular"}).Error; err != nil {
		return err
	}

	leaveTypes := []master.LeaveType{
		{Name: "Annual", DefaultQuota: 12, IsDeducted: true},
		{Name: "Sick", DefaultQuota: 15, IsDeducted: false},
		{Name: "Unpaid", DefaultQuota: 0, IsDeducted: false},
	}

	for _, lt := range leaveTypes {
		if err := tx.Where(master.LeaveType{Name: lt.Name}).FirstOrCreate(&lt).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedRolesAndAdmin(tx *gorm.DB, cfg *config.Config, hasher Hasher) (*rbac.Role, error) {
	roleSuperadmin := rbac.Role{Name: string(constants.UserRoleSuperadmin)}
	if err := tx.Where(rbac.Role{Name: roleSuperadmin.Name}).FirstOrCreate(&roleSuperadmin).Error; err != nil {
		return nil, err
	}

	roleEmployee := rbac.Role{Name: string(constants.UserRoleEmployee)}
	if err := tx.Where(rbac.Role{Name: roleEmployee.Name}).FirstOrCreate(&roleEmployee).Error; err != nil {
		return nil, err
	}

	newAdmin := user.User{
		Username: cfg.CredentialConfig.SuperadminUsername,
	}
	hashPass, err := hasher.HashPassword(cfg.CredentialConfig.SuperadminPassword)
	if err != nil {
		return nil, err
	}

	if err := tx.Where(user.User{Username: newAdmin.Username}).
		Attrs(user.User{
			PasswordHash:       hashPass,
			RoleID:             roleSuperadmin.ID,
			MustChangePassword: false,
			IsActive:           true,
		}).
		FirstOrCreate(&newAdmin).Error; err != nil {
		return nil, err
	}

	return &roleSuperadmin, nil
}

func seedPermissions(tx *gorm.DB, roleSuperadmin *rbac.Role) error {
	permissionsByGroup := []struct {
		GroupName   string
		Permissions []string
	}{
		{"Permission", []string{constants.VIEW_PERMISSION}},
		{"Role", []string{constants.CREATE_ROLE, constants.VIEW_ROLE, constants.ASSIGN_ROLE}},
		{"Master", []string{constants.VIEW_MASTER}},
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

	// Assign all permissions to superadmin role
	for _, pid := range permissionIDs {
		var rp rbac.RolePermission
		if err := tx.Where(rbac.RolePermission{RoleID: roleSuperadmin.ID, PermissionID: pid}).FirstOrCreate(&rp).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedCompanyData(tx *gorm.DB) error {
	// TODO: remove seed company after feature multi tenant implemented
	companyData := company.Company{Name: "PT. Pick", PhoneNumber: "08531432221023", Address: "Jl.Kejaksaan no.23 Jakarta Utara", Email: "admin@pick.com"}
	if err := tx.Where(company.Company{Name: companyData.Name}).FirstOrCreate(&companyData).Error; err != nil {
		return err
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
