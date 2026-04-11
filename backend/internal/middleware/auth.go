package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired is Gin middleware that checks for a valid JWT.
func AuthRequired() gin.HandlerFunc {
	// returns a function — Gin calls this function for every request to protected routes
	return func(ctx *gin.Context) {
		// read the Authorization header
		// expected format: "Bearer eyJhbGci..."
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
			// AbortWithStatusJSON stops the request here — the handler never runs
		}

		// split "Bearer eyJhbGci..." into ["Bearer", "eyJhbGci..."]
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		// validate the token (parts[1] is the actual JWT string)
		claims, err := ValidateToken(parts[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// extract user_id from claims and store it in the request context
		// handlers downstream can read this with ctx.GetInt("user_id")
		userID := int(claims["user_id"].(float64))
		// why float64? JSON numbers are always float64 in Go's encoding/json
		ctx.Set("user_id", userID)

		// pass the request to the next handler
		ctx.Next()
	}
}
