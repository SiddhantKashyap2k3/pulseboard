package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	// replace YOUR_USERNAME with your actual GitHub username
	"github.com/SiddhantKashyap2k3/pulseboard/internal/db"
)

func main() {

	// connect to the pulseboard database we just created in pgAdmin
	database, err := db.New(db.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "pulseboard",
	})
	if err != nil {
		// if DB connection fails, no point starting the server
		log.Fatal("Failed to connect to database:", err)
	}

	// when main() exits, cleanly close all DB connections
	defer database.Close()

	router := gin.Default()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "pulseboard-api",
		})
	})

	router.Run(":8080")
}