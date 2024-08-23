package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

// AppConfig holds the configuration values for the application
type AppConfig struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	ServerPort  string
	AutoMigrate bool
}

// Config is a global variable to access application configuration
var Config *AppConfig

// LoadConfig loads configuration from environment variables or .env file
func LoadConfig() *AppConfig {
	// Load .env file if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using environment variables.")
	}

	// Parse AUTO_MIGRATE environment variable
	autoMigrate, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE"))
	if err != nil {
		autoMigrate = false
	}
	// Create an AppConfig instance with environment variables
	Config = &AppConfig{
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		ServerPort:  os.Getenv("PORT"),
		AutoMigrate: autoMigrate,
	}

	return Config
}
