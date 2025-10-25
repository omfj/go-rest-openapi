package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/omfj/go-rest-openapi/docs"
)

// @title           Posts API
// @version         1.0
// @description     A REST API for managing posts and users
// @host            localhost:3000
// @BasePath        /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	state, err := NewServerFromEnv()
	if err != nil {
		log.Fatalf("Error initializing app state: %v", err)
	}

	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("GET /", state.handleHealthCheck)
	mux.HandleFunc("GET /posts", state.handleGetPosts)
	mux.HandleFunc("POST /posts", state.handleCreatePost)
	mux.HandleFunc("GET /user/{id}/posts", state.handleGetUserPosts)

	// API Documentation
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /scalar", scalarHandler)

	port := ":3000"
	fmt.Printf("🚀 Server starting on http://localhost%s\n", port)
	fmt.Printf("📚 Swagger UI @ http://localhost%s/swagger\n", port)
	fmt.Printf("📖 Scalar UI  @ http://localhost%s/scalar\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
