package seeder

import (
	"basekarya-backend/internal/config"
	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/logger"

	"gorm.io/gorm"
)

func Execute(db *gorm.DB, cfg *config.Config, hasher Hasher) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		generalDept := master.Department{Name: "Umum"}

		if err := tx.Where(master.Department{Name: "Umum"}).FirstOrCreate(&generalDept).Error; err != nil {
			return err
		}

		regularShift := master.Shift{Name: "Regular", StartTime: "09:00:00", EndTime: "18:00:00"}
		if err := tx.FirstOrCreate(&regularShift, master.Shift{Name: "Regular"}).Error; err != nil {
			return err
		}

		newAdmin := user.User{
			Username: cfg.CredentialConfig.SuperadminUsername,
		}
		hashPass, err := hasher.HashPassword(cfg.CredentialConfig.SuperadminPassword)
		if err != nil {
			return err
		}

		if err := tx.Where(user.User{Username: newAdmin.Username}).
			Attrs(user.User{
				PasswordHash:       hashPass,
				Role:               string(constants.UserRoleSuperadmin),
				MustChangePassword: false,
				IsActive:           true,
			}).
			FirstOrCreate(&newAdmin).Error; err != nil {
			return err
		}

		leaveTypeAnnual := master.LeaveType{Name: "Annual", DefaultQuota: 12, IsDeducted: true}
		leaveTypeSick := master.LeaveType{Name: "Sick", DefaultQuota: 15, IsDeducted: false}
		leaveTypeUnpaid := master.LeaveType{Name: "Unpaid", DefaultQuota: 0, IsDeducted: false}

		if err := tx.Where(master.LeaveType{Name: leaveTypeAnnual.Name}).FirstOrCreate(&leaveTypeAnnual).Error; err != nil {
			return err
		}

		if err := tx.Where(master.LeaveType{Name: leaveTypeSick.Name}).FirstOrCreate(&leaveTypeSick).Error; err != nil {
			return err
		}

		if err := tx.Where(master.LeaveType{Name: leaveTypeUnpaid.Name}).FirstOrCreate(&leaveTypeUnpaid).Error; err != nil {
			return err
		}

		companyData := company.Company{Name: "PT. Pick", PhoneNumber: "08531432221023", Address: "Jl.Kejaksaan no.23 Jakarta Utara", Email: "admin@pick.com"}

		if err := tx.Where(company.Company{Name: companyData.Name}).FirstOrCreate(&companyData).Error; err != nil {
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
