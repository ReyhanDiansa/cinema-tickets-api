package middleware

import (
	"cinema-tickets/utils"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)


// JWT Claims structure
type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

// AuthMiddleware validates the JWT and ensures that the user is authenticated
func AuthMiddleware() gin.HandlerFunc {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			utils.ResponseFormatter(c, 401, false, "Authorization token required", nil)
			c.Abort()
			return
		}

		// The token comes in the format "Bearer <token>"
		parts := strings.Split(tokenStr, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ResponseFormatter(c, 401, false, "Invalid token format", nil)
			c.Abort()
			return
		}

		tokenStr = parts[1] // Extract the token from the Bearer string

		// Parse the token
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Return the secret key used to sign the token
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			utils.ResponseFormatter(c, 401, false, "Invalid or expired token", nil)
			c.Abort()
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(*Claims)
		if !ok {
			utils.ResponseFormatter(c, 401, false, "Invalid token claims", nil)
			c.Abort()
			return
		}

		// Store the claims in the context to use in further routes
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware checks the user's role and ensures they have access to the resource
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the role from the context (set by AuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			utils.ResponseFormatter(c, 401, false, "Role not found in token", nil)			
			c.Abort()
			return
		}

		// Check if the user's role matches one of the required roles
		authorized := false
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				authorized = true
				break
			}
		}

		if !authorized {
			utils.ResponseFormatter(c, 403, false, "You are not authorized to access this resource", nil)						
			c.Abort()
			return
		}

		c.Next()
	}
}
