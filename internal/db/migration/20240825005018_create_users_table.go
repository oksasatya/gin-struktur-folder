package migration

import (
	"gin-struktur-folder/internal/app/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Fungsi migrasi untuk create_users_table
func createUsersTableMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240825005018_create_users_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&model.User{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("users")
		},
	}
}
