package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/SiddhantKashyap2k3/pulseboard/internal/models"
)

// generateAPIKey creates a cryptographically secure random 64-char hex string.
func generateAPIKey() (string, error) {
	// 32 bytes of random data from the OS entropy source
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate API key: %w", err)
	}
	// Encode to hex: 32 bytes → 64 hex characters
	return hex.EncodeToString(bytes), nil
}

// CreateWorkspace inserts a new workspace with an auto-generated API key.
func CreateWorkspace(db *sql.DB, name string, userID int) (*models.Workspace, error) {
	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, err
	}

	ws := &models.Workspace{}
	err = db.QueryRow(`
		INSERT INTO workspaces (name, api_key, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, api_key, user_id, created_at`,
		name, apiKey, userID,
	).Scan(&ws.ID, &ws.Name, &ws.APIKey, &ws.UserID, &ws.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating workspace: %w", err)
	}
	return ws, nil
}

// ListWorkspacesByUser returns all workspaces owned by a specific user.
func ListWorkspacesByUser(db *sql.DB, userID int) ([]models.Workspace, error) {
	rows, err := db.Query(`
		SELECT id, name, api_key, user_id, created_at
		FROM workspaces WHERE user_id = $1
		ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error listing workspaces: %w", err)
	}
	// defer rows.Close() ensures the DB connection is returned to the pool
	// even if we return early due to an error during row scanning.
	// Without this, the connection leaks — the pool shrinks over time
	// and eventually your app runs out of connections.
	defer rows.Close()

	var workspaces []models.Workspace
	for rows.Next() {
		var ws models.Workspace
		if err := rows.Scan(&ws.ID, &ws.Name, &ws.APIKey, &ws.UserID, &ws.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning workspace: %w", err)
		}
		workspaces = append(workspaces, ws)
	}
	return workspaces, nil
}

// GetWorkspaceByAPIKey finds a workspace by its API key.
// Used by the API key middleware to authenticate incoming events.
func GetWorkspaceByAPIKey(db *sql.DB, apiKey string) (*models.Workspace, error) {
	ws := &models.Workspace{}
	err := db.QueryRow(`
		SELECT id, name, api_key, user_id, created_at
		FROM workspaces WHERE api_key = $1`, apiKey,
	).Scan(&ws.ID, &ws.Name, &ws.APIKey, &ws.UserID, &ws.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid API key")
		}
		return nil, fmt.Errorf("error fetching workspace: %w", err)
	}
	return ws, nil
}