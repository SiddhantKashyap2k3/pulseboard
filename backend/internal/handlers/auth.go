package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/SiddhantKashyap2k3/pulseboard/internal/db"
	"github.com/SiddhantKashyap2k3/pulseboard/internal/middleware"
	"github.com/SiddhantKashyap2k3/pulseboard/internal/models"
)

// AuthHandler holds the database connection
// all handler methods hang off this struct
// this is called "dependency injection" — we pass the DB in
// instead of using a global variable
type AuthHandler struct {
	DB *sql.DB
}

// Register handles POST /register
func (h *AuthHandler) Register(ctx *gin.Context) {
	// ShouldBindJSON reads the request body JSON into our struct
	// it also runs the binding validations we defined (required, email, min=6)
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 400 Bad Request — client sent invalid data
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// bcrypt.GenerateFromPassword hashes the password
	// bcrypt.DefaultCost (10) controls how slow the hashing is
	// higher cost = slower = harder to brute force
	// we convert the string to []byte because bcrypt works on bytes
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		// 500 Internal Server Error — something went wrong on our side
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// insert the user into the database
	user, err := db.CreateUser(h.DB, req.Email, string(hashedPassword))
	if err != nil {
		// if email already exists Postgres returns a unique constraint error
		ctx.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	// 201 Created — resource was successfully created
	// we return only safe fields — never the password hash
	ctx.JSON(http.StatusCreated, models.RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find the user by email
	user, err := db.GetUserByEmail(h.DB, req.Email)
	if err != nil {
		// don't reveal whether the email exists or not — security best practice
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// bcrypt.CompareHashAndPassword takes the stored hash and the plain password
	// it hashes the plain password with the same salt and compares
	// returns nil if they match, error if they don't
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// password matched — generate a JWT
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, models.LoginResponse{Token: token})
}
