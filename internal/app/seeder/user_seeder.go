package seeder

import (
	"gin-struktur-folder/internal/app/model"
	"gin-struktur-folder/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	var users []model.User

	for i := 1; i <= 15; i++ {
		password, _ := utils.HashPassword("test12345")

		user := model.User{
			Email:     utils.RandomEmail(),
			FirstName: utils.RandomFirstName(),
			LastName:  utils.RandomLastName(),
			Password:  password,
		}
		users = append(users, user)
	}
	if err := db.Create(&users).Error; err != nil {
		return
	}
	logrus.Println("Seed users success")
}
