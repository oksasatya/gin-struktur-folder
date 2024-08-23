package model

import (
	"gorm.io/gorm"
)

// User represents the user model for the database
type User struct {
	gorm.Model
	FirstName string `gorm:"size:255" json:"first_name" binding:"required,min=2,max=255"`
	LastName  string `gorm:"size:255" json:"last_name" binding:"required,min=2,max=255"`
	Email     string `gorm:"unique;not null;size:255" json:"email" binding:"required,email"`
	Password  string `gorm:"not null" json:"password" binding:"required,min=6"`
}

// LoginUser represents the user model for login
type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// TableName returns the table name for the user model
func (User) TableName() string {
	return "users"
}
