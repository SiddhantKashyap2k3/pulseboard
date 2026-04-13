package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SiddhantKashyap2k3/pulseboard/internal/db"
	"github.com/SiddhantKashyap2k3/pulseboard/internal/models"
)

// WorkspaceHandler holds the database connection for workspace routes.
type WorkspaceHandler struct {
	DB *sql.DB
}

// Create handles POST /api/v1/workspaces
func (h *WorkspaceHandler) Create(ctx *gin.Context) {
	// user_id was set by AuthRequired middleware from the JWT
	userID := ctx.GetInt("user_id")

	var req models.CreateWorkspaceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ws, err := db.CreateWorkspace(h.DB, req.Name, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create workspace"})
		return
	}

	ctx.JSON(http.StatusCreated, models.WorkspaceResponse{
		ID:        ws.ID,
		Name:      ws.Name,
		APIKey:    ws.APIKey,
		CreatedAt: ws.CreatedAt,
	})
}

// List handles GET /api/v1/workspaces
func (h *WorkspaceHandler) List(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")

	workspaces, err := db.ListWorkspacesByUser(h.DB, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list workspaces"})
		return
	}

	// Build response slice — don't expose internal fields like UserID
	resp := make([]models.WorkspaceResponse, len(workspaces))
	for i, ws := range workspaces {
		resp[i] = models.WorkspaceResponse{
			ID:        ws.ID,
			Name:      ws.Name,
			APIKey:    ws.APIKey,
			CreatedAt: ws.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, resp)
}