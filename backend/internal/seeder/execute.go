package seeder

import (
	"hris-backend/internal/config"
	"hris-backend/internal/modules/master"
	"hris-backend/internal/modules/user"
	"hris-backend/pkg/constants"
	"hris-backend/pkg/logger"

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

		employeesToSeed := []struct {
			NIK      string
			Name     string
			DeptID   uint
			ShiftID  uint
			Role     string
			IsActive bool
		}{
			{"EMP001", "Taufik Januar", generalDept.ID, regularShift.ID, string(constants.UserRoleEmployee), true},
			//TODO: insert real employees later
		}

		for _, empData := range employeesToSeed {
			hashPass, _ := hasher.HashPassword(empData.NIK)
			newUser := user.User{
				Username:           empData.NIK,
				PasswordHash:       hashPass,
				Role:               empData.Role,
				MustChangePassword: true,
				IsActive:           empData.IsActive,
			}

			if err := tx.Create(&newUser).Error; err != nil {
				logger.Errorf("User %s already exists or error: %v", empData.NIK, err)
				continue
			}
			newEmployee := user.Employee{
				UserID:       newUser.ID,
				DepartmentID: empData.DeptID,
				ShiftID:      empData.ShiftID,
				NIK:          empData.NIK,
				FullName:     empData.Name,
			}

			if err := tx.Create(&newEmployee).Error; err != nil {
				return err
			}

			logger.Infof("Seeded: %s - %s", empData.NIK, empData.Name)
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
