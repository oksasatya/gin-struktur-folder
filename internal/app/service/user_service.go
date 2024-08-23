package service

import (
	"errors"
	"gin-struktur-folder/internal/app/model"
	"gin-struktur-folder/internal/app/repository"
	"gin-struktur-folder/internal/middleware"
	"gin-struktur-folder/pkg/utils"
)

type UserService interface {
	Register(user *model.User) error
	Login(loginUser *model.LoginUser) (string, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository, jwtSecret string) UserService {
	return &userService{repo: repo, jwtSecret: jwtSecret}
}

func (s *userService) Register(user *model.User) error {
	// check if user already exists
	if _, err := s.repo.GetUserByEmail(user.Email); err == nil {
		return errors.New("email already exists")
	}

	// hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// create the user
	return s.repo.CreateUser(user)
}

func (s *userService) Login(loginUser *model.LoginUser) (string, error) {
	user, err := s.repo.GetUserByEmail(loginUser.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// check if the password is correct
	if !utils.CheckPasswordHash(loginUser.Password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	//generate jwt token
	token, err := middleware.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
