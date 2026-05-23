package seeder

import (
	"strings"

	"basekarya-backend/internal/config"
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/onboarding"
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

		if err := seedOnboardingTemplates(tx); err != nil {
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
	generalDept := master.Department{Name: "Umum", CompanyID: 1}
	if err := tx.Where(master.Department{Name: "Umum"}).FirstOrCreate(&generalDept).Error; err != nil {
		return err
	}

	regularShift := master.Shift{Name: "Regular", StartTime: "09:00:00", EndTime: "18:00:00", CompanyID: 1}
	if err := tx.FirstOrCreate(&regularShift, master.Shift{Name: "Regular"}).Error; err != nil {
		return err
	}

	leaveTypes := []master.LeaveType{
		{Name: "Annual", DefaultQuota: 12, IsDeducted: true, CompanyID: 1},
		{Name: "Sick", DefaultQuota: 15, IsDeducted: false, CompanyID: 1},
		{Name: "Unpaid", DefaultQuota: 0, IsDeducted: false, CompanyID: 1},
	}

	for _, lt := range leaveTypes {
		if err := tx.Where(master.LeaveType{Name: lt.Name}).FirstOrCreate(&lt).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedRolesAndAdmin(tx *gorm.DB, cfg *config.Config, hasher Hasher) (*rbac.Role, error) {
	platformRole := rbac.Role{Name: string(constants.UserRolePlatformAdmin), CompanyID: 0}
	if err := tx.Where(rbac.Role{Name: string(constants.UserRolePlatformAdmin)}).FirstOrCreate(&platformRole).Error; err != nil {
		return nil, err
	}

	newAdmin := user.User{
		Username:        cfg.CredentialConfig.SuperadminUsername,
		CompanyID:       0,
		IsPlatformAdmin: true,
	}
	hashPass, err := hasher.HashPassword(cfg.CredentialConfig.SuperadminPassword)
	if err != nil {
		return nil, err
	}

	if err := tx.Where(user.User{Username: newAdmin.Username}).
		Attrs(user.User{
			PasswordHash:       hashPass,
			RoleID:             platformRole.ID,
			CompanyID:          0,
			IsPlatformAdmin:    true,
			MustChangePassword: false,
			IsActive:           true,
		}).
		FirstOrCreate(&newAdmin).Error; err != nil {
		return nil, err
	}

	if newAdmin.CompanyID != 0 {
		newAdmin.CompanyID = 0
		newAdmin.RoleID = platformRole.ID
		newAdmin.IsPlatformAdmin = true
		if err := tx.Save(&newAdmin).Error; err != nil {
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
		{"Onboarding", []string{constants.VIEW_ONBOARDING, constants.MANAGE_ONBOARDING_TEMPLATE, constants.UPDATE_ONBOARDING_TASK}},
		{"Finance", []string{constants.VIEW_FINANCE, constants.CREATE_FINANCE, constants.APPROVAL_FINANCE, constants.EXPORT_FINANCE, constants.MANAGE_FINANCE_CATEGORY, constants.VIEW_FINANCE_DASHBOARD}},
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
		{Name: "Pro", Slug: "pro", MaxEmployees: 0, PriceMonthly: 249000, Features: `{"modules":["attendance","leave","overtime","loan","reimbursement","payroll","contract","finance","recruitment","onboarding"]}`, IsActive: true},
	}

	for i := range plans {
		p := &plans[i]
		var existing subscription.SubscriptionPlan
		err := tx.Where("slug = ?", p.Slug).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := tx.Create(p).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
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

func seedOnboardingTemplates(tx *gorm.DB) error {
	type templateDef struct {
		Name       string
		Department string
		Tasks      []struct {
			Name, Description string
			Order             int
		}
	}

	templates := []templateDef{
		{
			Name:       "IT Setup",
			Department: "IT",
			Tasks: []struct {
				Name, Description string
				Order             int
			}{
				{"Buat akun email perusahaan", "Buat akun email @company.com untuk karyawan baru", 1},
				{"Setup laptop/PC", "Siapkan laptop atau PC beserta perangkat yang dibutuhkan", 2},
				{"Berikan akses ke sistem internal", "Berikan akses ke aplikasi HR, project management, dan sistem internal lainnya", 3},
				{"Instalasi software wajib", "Install software wajib sesuai kebutuhan divisi", 4},
			},
		},
		{
			Name:       "HR Document Collection",
			Department: "HR",
			Tasks: []struct {
				Name, Description string
				Order             int
			}{
				{"Kumpulkan KTP", "Minta dan simpan salinan KTP karyawan baru", 1},
				{"Kumpulkan Kartu Keluarga (KK)", "Minta dan simpan salinan Kartu Keluarga", 2},
				{"Kumpulkan NPWP", "Minta dan simpan salinan NPWP", 3},
				{"Kumpulkan BPJS Kesehatan & Ketenagakerjaan", "Minta dan simpan kartu BPJS atau nomor kepesertaan", 4},
				{"Upload foto pas 3x4", "Minta foto resmi berukuran 3x4 untuk identitas karyawan", 5},
				{"Tanda tangan kontrak kerja", "Proses penandatanganan kontrak kerja", 6},
			},
		},
	}

	for _, def := range templates {
		var existing onboarding.OnboardingTemplate
		result := tx.Where("name = ? AND department = ?", def.Name, def.Department).First(&existing)
		if result.Error == nil {
			continue
		}

		tmpl := onboarding.OnboardingTemplate{
			Name:       def.Name,
			Department: def.Department,
			CompanyID:  1,
		}
		for _, task := range def.Tasks {
			tmpl.Items = append(tmpl.Items, onboarding.OnboardingTemplateItem{
				TaskName:    task.Name,
				Description: task.Description,
				SortOrder:   task.Order,
				CompanyID:   1,
			})
		}
		if err := tx.Create(&tmpl).Error; err != nil {
			return err
		}
	}

	return nil
}
