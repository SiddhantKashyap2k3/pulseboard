package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Config holds all the values needed to connect to Postgres
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// New creates a connection pool and returns it
func New(cfg Config) (*sql.DB, error) {

	// build the connection string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	// sql.Open validates the driver and DSN format but does NOT connect yet
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// db.Ping() actually attempts a real connection
	// if Postgres is not running or password is wrong this fails here
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	// max 25 connections open at the same time
	db.SetMaxOpenConns(25)

	// 5 connections stay open and ready even when idle
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db, nil
}