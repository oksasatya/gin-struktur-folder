package controller

import (
	"errors"
	"gin-struktur-folder/internal/app/model"
	"gin-struktur-folder/internal/app/service"
	"gin-struktur-folder/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

func (uc *UserController) Register(c *gin.Context) {
	var user model.User
	// bind the request body to the user model
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.HandleError(c, err, "Invalid Request")
		return
	}

	// register the user
	if err := uc.service.Register(&user); err != nil {
		utils.HandleError(c, err, "Failed to register user")
		return
	}
	// respond with the user
	//utils.RespondJSON(c, http.StatusCreated, user, "User Register successfully")
	utils.RespondJSON(c, http.StatusCreated, map[string]any{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}, "User created successfully")
}

func (uc *UserController) Login(c *gin.Context) {
	var loginUser model.LoginUser
	// bind the request body to the login user model
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		utils.HandleError(c, err, "Invalid Request")
		return
	}
	// validate the login user
	if err := validator.New().Struct(loginUser); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			formatter := utils.FormatValidationError(loginUser, validationErrors)
			utils.RespondJSON(c, http.StatusBadRequest, nil, formatter)
			return
		}
		utils.RespondJSON(c, http.StatusBadRequest, nil, "Invalid Request : "+err.Error())
		return
	}

	// Auth the user
	token, err := uc.service.Login(&loginUser)
	if err != nil {
		utils.RespondJSON(c, http.StatusUnauthorized, nil, "Invalid email or password")
		return
	}

	// respond with the user
	utils.RespondJSON(c, http.StatusOK, gin.H{"token": token}, "User logged in successfully")
}
