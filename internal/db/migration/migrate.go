package migration

import (
	"fmt"
	"gin-struktur-folder/internal/app/seeder"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
)

// Migrate is a function to migrate all tables
func Migrate(db *gorm.DB) error {
	migrations := []*gormigrate.Migration{
		createUsersTableMigration(),
	}

	m := gormigrate.New(db, &gormigrate.Options{
		TableName:                 "migrations",
		IDColumnName:              "id",
		IDColumnSize:              255,
		UseTransaction:            true,
		ValidateUnknownMigrations: true,
	}, migrations)

	autoMigrate := os.Getenv("AUTO_MIGRATE")
	autoDrop := os.Getenv("AUTO_DROP")

	if autoDrop == "true" && autoMigrate == "true" {
		logrus.Println("Running AutoDrop (Rollback all migrations) and AutoMigrate...")

		for i := len(migrations) - 1; i >= 0; i-- {
			if err := m.RollbackMigration(migrations[i]); err != nil {
				return fmt.Errorf("rollback migration %s failed: %v", migrations[i].ID, err)
			}
		}

		if err := m.Migrate(); err != nil {
			return fmt.Errorf("migration failed after drop: %v", err)
		}
		logrus.Println("Running Seeders...")
		seedAll(db)
		logrus.Println("AutoMigrate and Seeders completed.")
	} else if autoDrop == "true" {
		logrus.Println("Running AutoDrop (Rollback all migrations)...")

		for i := len(migrations) - 1; i >= 0; i-- {
			if err := m.RollbackMigration(migrations[i]); err != nil {
				return fmt.Errorf("rollback migration %s failed: %v", migrations[i].ID, err)
			}
		}
	} else if autoMigrate == "true" {
		logrus.Println("Running AutoMigrate...")
		if err := m.Migrate(); err != nil {
			return fmt.Errorf("migration failed: %v", err)
		}

		logrus.Println("Running Seeders...")
		seedAll(db)
		logrus.Println("AutoMigrate and Seeders completed.")
	} else {
		logrus.Println("Skipping AutoMigrate and AutoDrop.")
	}

	return nil
}

func seedAll(db *gorm.DB) {
	seeder.SeedUsers(db)
	logrus.Println("Seed all success")
}
