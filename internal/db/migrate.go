package db

import (
	"gin-struktur-folder/internal/app/model"
	"gin-struktur-folder/internal/app/seeder"
	"github.com/sirupsen/logrus"
)

// AutoMigrate migrates all models to the database
func AutoMigrate() {
	err := DB.AutoMigrate(
		&model.User{},
	)
	if err != nil {
		logrus.Fatalf("Failed to migrate database: %v", err)
	}
	logrus.Println("Database migration completed")
}

// AutoDrop drops all tables in the database
func AutoDrop() {
	err := DB.Migrator().DropTable(
		&model.User{},
	)
	if err != nil {
		logrus.Fatalf("Failed to drop database tables: %v", err)
	}
	logrus.Println("Database tables dropped")
}

// SeedAll seeds the database with initial data
func SeedAll() {
	seeder.SeedUsers(DB)
}
