package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SiddhantKashyap2k3/pulseboard/internal/db"
	"github.com/SiddhantKashyap2k3/pulseboard/internal/handlers"
	"github.com/SiddhantKashyap2k3/pulseboard/internal/middleware"
)

func main() {
	// connect to database
	database, err := db.New(db.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "pulseboard",
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// create handler instances — passing the DB connection in
	authHandler := &handlers.AuthHandler{DB: database}

	// create workspce handler instance - for DB conncetion not crete own DB connection but from here --> dependecy injection
	workspaceHandler := &handlers.WorkspaceHandler{DB: database}

	router := gin.Default()

	// health check route
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "pulseboard-api",
		})
	})

	// auth routes grouped under /api/v1
	// grouping keeps URLs organised and lets us version our API
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register) // POST /api/v1/auth/register
			auth.POST("/login", authHandler.Login)       // New - Login route
		}
	}

	// protected routes — JWT required
	// middleware.AuthRequired() runs BEFORE every handler in this group
	protected := v1.Group("")
	protected.Use(middleware.AuthRequired())
	{
		// test route — returns the logged-in user's ID
		// proves the middleware extracted user_id from the token
		protected.GET("/me", func(ctx *gin.Context) {
			userID := ctx.GetInt("user_id")
			ctx.JSON(http.StatusOK, gin.H{"user_id": userID})
		})

		// Tell GIN when request hits workspaces run first create then list
		// Flow --> Request-> AuthRequired middleware -> checks JWT -> sets user_id -> handler runs
		protected.POST("/workspaces", workspaceHandler.Create)
		protected.GET("/workspaces", workspaceHandler.List)
	}

	router.Run(":8080")
}
