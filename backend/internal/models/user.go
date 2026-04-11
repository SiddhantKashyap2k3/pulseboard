package models

import "time"

// User mirrors the users table in the database
// Each field maps to one column
type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"` // json:"-" means this field is NEVER included
	// in JSON responses — password never leaves server
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterRequest is what the client sends in the request body
// when calling POST /register
type RegisterRequest struct {
	Email string `json:"email" binding:"required,email"` // binding:"required" means
	// Gin returns 400 if missing
	// binding:"email" validates format
	Password string `json:"password" binding:"required,min=6"` // min=6 means minimum 6 chars
}

// RegisterResponse is what we send back after successful registration
// we never send the password hash back
type RegisterResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest is the body for POST /login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is returned after successful login.
type LoginResponse struct {
	Token string `json:"token"`
}
