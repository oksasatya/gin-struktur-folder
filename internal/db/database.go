package db

import (
	"fmt"
	"gin-struktur-folder/internal/config"
	"gin-struktur-folder/internal/db/migration"
	"gin-struktur-folder/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.Config
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	newLogger := utils.NewLogrusGormLogger(logger.Info)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		logrus.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := migration.Migrate(DB); err != nil {
		logrus.Fatalf("Failed to migrate: %v", err)
	}
}
