package main

import (
	"net/http"
)

type User struct {
	ID        int
	Username  string
	Email     string
	CreatedAt string
}

type Session struct {
	ID           int
	UserID       int
	SessionToken string
	ExpiresAt    string
}

type Auth struct {
	user    User
	session Session
}

func (s *Server) GetAuthFromRequest(r *http.Request) (*Auth, error) {
	var auth Auth

	// Get beearer token from request
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, nil // No auth provided
	}

	// Trim Bearer prefix
	token = token[len("Bearer "):]

	// Get session with token
	var session Session
	var query = `SELECT id, user_id, session_token, expires_at FROM sessions WHERE session_token = ?`

	err := s.Pool.QueryRow(query, token).Scan(&session.ID, &session.UserID, &session.SessionToken, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	// Get user with session's user_id
	var user User
	query = `SELECT id, username, email, created_at FROM users WHERE id = ?`
	err = s.Pool.QueryRow(query, session.UserID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	auth.user = user
	auth.session = session

	return &auth, nil
}
