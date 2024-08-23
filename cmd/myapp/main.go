package main

import (
	"context"
	"errors"
	"gin-struktur-folder/internal/app/routes"
	"gin-struktur-folder/internal/config"
	"gin-struktur-folder/internal/db"
	"gin-struktur-folder/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// global variable to access resty client
var client *resty.Client

func main() {
	// load the configuration
	config.LoadConfig()
	// Initialize the logger
	utils.SetupLogger()
	// init DB
	db.InitDB()
	// init gin router
	r := gin.New()
	// Add middleware for logging requests
	r.Use(utils.LogrusLogger())
	// init controller
	c := routes.InitController(db.DB, config.Config.JWTSecret)
	// register routes
	routes.RegisterRoutes(r, c, db.DB, []byte(config.Config.JWTSecret))
	// init resty client
	client = resty.New()
	client.SetTimeout(5 * time.Second)

	// start the server
	startServer(r)
}

// startServer is a function to start the server
func startServer(r *gin.Engine) {
	srv := &http.Server{
		Addr:    ":" + config.Config.ServerPort,
		Handler: r,
	}

	// run server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()
	logrus.Println("Server is running on port:", config.Config.ServerPort)

	// wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds
	quit := make(chan os.Signal, 1)
	// catch SIGINT (Ctrl+C) and SIGTERM (Docker container stop)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Println("Shutting down server...")

	// create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Println("Server exiting")
}
