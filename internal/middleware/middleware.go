package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

// JWTMiddleware is a middleware for validating JWT tokens in Gin
func JWTMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or malformed JWT"})
			c.Abort()
			return
		}

		// Check if the token is in the format "Bearer {token}"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			var status int
			var message string

			switch err.Error() {
			case jwt.ErrSignatureInvalid.Error():
				status = http.StatusUnauthorized
				message = "Invalid signature"
			case "token is expired":
				status = http.StatusUnauthorized
				message = "Token is expired, please login again"
			default:
				status = http.StatusUnauthorized
				message = "Unauthorized access"
			}

			c.JSON(status, gin.H{"error": message})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired JWT token"})
			c.Abort()
			return
		}

		// Store user information in context if needed
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			c.Set("user", claims)
		}

		// Continue to the next handler
		c.Next()
	}
}

// GenerateToken generates a JWT token
func GenerateToken(userId uint, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// RoleMiddleware is a middleware for checking user roles in Gin
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userClaims, ok := claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		role, ok := userClaims["role"].(string)
		if !ok || role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this route"})
			c.Abort()
			return
		}

		// Continue to the next handler
		c.Next()
	}
}
