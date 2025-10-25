package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type Server struct {
	Pool *sql.DB
}

func runMigrations(db *sql.DB) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, "./migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func NewServerFromEnv() (*Server, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return &Server{Pool: db}, nil
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
