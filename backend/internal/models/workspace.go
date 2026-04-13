package models

import "time"

// Workspace mirrors the workspaces table in the database.
type Workspace struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateWorkspaceRequest is the body for POST /api/v1/workspaces.
type CreateWorkspaceRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}

// WorkspaceResponse is returned after creating or listing workspaces.
// We return the API key on creation, but you could choose to mask it on list.
type WorkspaceResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}
