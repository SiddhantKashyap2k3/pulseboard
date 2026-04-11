package db

import (
	"database/sql"
	"fmt"

	"github.com/SiddhantKashyap2k3/pulseboard/internal/models"
)

// CreateUser inserts a new user into the database
// It takes the db connection pool, email, and hashed password
// It returns the created user or an error
func CreateUser(db *sql.DB, email, passwordHash string) (*models.User, error) {
	// user is a pointer to a models.User struct — we'll fill it with data from Postgres
	user := &models.User{}

	// QueryRow executes a SQL query that returns exactly one row
	// $1, $2, $3 are placeholders — Go fills them in safely
	// using placeholders instead of string concatenation prevents SQL injection attacks
	// RETURNING tells Postgres to send back the values it just inserted
	err := db.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, created_at, updated_at`,
		email,        // replaces $1
		passwordHash, // replaces $2
	).Scan(
		// Scan reads the returned columns into our struct fields
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

// GetUserByEmail finds a user by their email address
// Used during login to look up the user before checking password
func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	user := &models.User{}

	err := db.QueryRow(`
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		// sql.ErrNoRows means no user found with that email
		// we return a clear error message instead of a confusing internal one
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return user, nil
}
