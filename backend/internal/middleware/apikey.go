package middleware

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SiddhantKashyap2k3/pulseboard/internal/db"
)

// APIKeyAuth is Gin middleware that authenticates requests via the X-API-Key header.
// It looks up the workspace by API key and sets workspace_id in the context.
func APIKeyAuth(database *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-Key")
		if apiKey == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing X-API-Key header"})
			return
		}

		ws, err := db.GetWorkspaceByAPIKey(database, apiKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			return
		}

		// Store workspace_id so downstream handlers know which tenant this request belongs to
		ctx.Set("workspace_id", ws.ID)
		ctx.Next()
	}
}
