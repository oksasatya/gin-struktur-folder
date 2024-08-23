package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

// RespondJSON is a utility function that sends a JSON response to the client using Gin
func RespondJSON(c *gin.Context, status int, payload interface{}, message string) {
	response := map[string]interface{}{
		"message": message,
		"data":    payload,
	}
	c.JSON(status, response)
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash checks if a password matches its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// FormatValidationError formats the validation error messages to a more readable form
func FormatValidationError(model interface{}, errors validator.ValidationErrors) string {
	if model == nil {
		// Handle the case where model is nil
		return "Invalid input: model cannot be nil"
	}
	var errorMessages []string
	modelType := reflect.TypeOf(model).Elem()

	for _, err := range errors {
		// Get the field name from the model's struct
		field, found := modelType.FieldByName(err.Field())
		fieldName := err.Field()
		if found {
			fieldName = field.Tag.Get("json")
			if fieldName == "" {
				fieldName = err.Field()
			}
		}

		// Generate custom error message based on the validation tag
		var message string
		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", fieldName)
		case "numeric":
			message = fmt.Sprintf("%s must be a number", fieldName)
		case "unique":
			message = fmt.Sprintf("%s already exists", fieldName)
		case "email":
			message = fmt.Sprintf("%s must be a valid email address", fieldName)
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters", fieldName, err.Param())
		case "max":
			message = fmt.Sprintf("%s must not be longer than %s characters", fieldName, err.Param())
		case "gte":
			message = fmt.Sprintf("%s must be greater than or equal to %s", fieldName, err.Param())
		case "gt":
			message = fmt.Sprintf("%s must be greater than %s", fieldName, err.Param())
		case "lte":
			message = fmt.Sprintf("%s must be less than or equal to %s", fieldName, err.Param())
		case "lt":
			message = fmt.Sprintf("%s must be less than %s", fieldName, err.Param())
		default:
			message = fmt.Sprintf("%s is invalid", fieldName)
		}

		errorMessages = append(errorMessages, message)
	}

	return strings.Join(errorMessages, ", ")
}

// IsValidUrl checks if a string is a valid URL
func IsValidUrl(url string) bool {
	var urlRegex = `^(http[s]?:\/\/)?[^\s(["<,>]*\.[^\s[",><]*$`
	return regexp.MustCompile(urlRegex).MatchString(url)
}

// HandleError handles errors and responds with appropriate messages
func HandleError(c *gin.Context, err error, message string) {
	// Check if it's a validation error
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errMsg := FormatValidationError(c, validationErrors)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errMsg,
			"message": "Invalid Request",
		})
		return
	}

	// For other errors, respond with a generic message
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   message,
		"message": err.Error(),
	})
}
