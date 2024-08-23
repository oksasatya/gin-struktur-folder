package routes

import (
	"gin-struktur-folder/internal/app/controller"
	"gin-struktur-folder/internal/app/repository"
	"gin-struktur-folder/internal/app/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewController struct {
	UserController controller.UserController
}

// RegisterRoutes is a function to register all routes
func RegisterRoutes(e *gin.Engine, c *NewController, db *gorm.DB, jwtSecret []byte) {
	// User routes
	r := e.Group("users")
	r.POST("/register", c.UserController.Register)
	r.POST("/login", c.UserController.Login)

	// Protected routes
	//protectedRoutes := e.Group("/user")
	//protectedRoutes.Use(middleware.JWTMiddleware(jwtSecret))
	//{
	//}
}

func InitController(db *gorm.DB, jwtSecret string) *NewController {
	return &NewController{
		UserController: controller.NewUserController(service.NewUserService(repository.NewUserRepository(db), jwtSecret)),
	}
}
