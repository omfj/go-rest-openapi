package main

import (
	"encoding/json"
	"net/http"
)

type Post struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// handleHealthCheck returns OK if the server is running
// @Summary      Health check
// @Description  Returns OK if the server is running
// @Tags         health
// @Produce      plain
// @Success      200  {string}  string  "OK"
// @Router       / [get]
func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// handleGetPosts returns all posts from the database
// @Summary      Get all posts
// @Description  Returns all posts ordered by creation date (newest first)
// @Tags         posts
// @Produce      json
// @Success      200  {array}   Post
// @Failure      500  {string}  string  "Failed to fetch posts"
// @Router       /posts [get]
func (s *Server) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []Post
	var query = `SELECT id, user_id, title, content, created_at FROM posts ORDER BY created_at DESC`

	rows, err := s.Pool.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	JSON(w, http.StatusOK, posts)
}

// handleGetUserPosts returns all posts for a specific user
// @Summary      Get user posts
// @Description  Returns all posts for a specific user ordered by creation date (newest first)
// @Tags         posts
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {array}   Post
// @Failure      500  {string}  string  "Failed to fetch user posts"
// @Router       /user/{id}/posts [get]
func (s *Server) handleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	var posts []Post
	var query = `SELECT id, user_id, title, content, created_at FROM posts WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := s.Pool.Query(query, userID)
	if err != nil {
		http.Error(w, "Failed to fetch user posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	JSON(w, http.StatusOK, posts)
}

type CreatePostBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// handleCreatePost creates a new post
// @Summary      Create post
// @Description  Creates a new post for the authenticated user
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        post  body      CreatePostBody  true  "Post data"
// @Success      201   {object}  Post
// @Failure      400   {string}  string  "Invalid request body"
// @Failure      401   {string}  string  "Failed to authenticate"
// @Failure      500   {string}  string  "Failed to create post"
// @Router       /posts [post]
func (s *Server) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	auth, err := s.GetAuthFromRequest(r)
	if err != nil {
		http.Error(w, "Failed to authenticate", http.StatusUnauthorized)
		return
	}

	var req CreatePostBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var post Post
	query := `INSERT INTO posts (user_id, title, content, created_at)
	          VALUES (?, ?, ?, datetime('now'))
	          RETURNING id, user_id, title, content, created_at`

	err = s.Pool.QueryRow(query, auth.user.ID, req.Title, req.Content).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt,
	)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusCreated, post)
}
