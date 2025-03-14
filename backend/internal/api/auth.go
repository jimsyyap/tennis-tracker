package api

import (
	"encoding/json"
	"net/http"

	"github.com/jimsyyap/tennis-tracker/backend/internal/database"
	"github.com/jimsyyap/tennis-tracker/backend/internal/middleware"
	"github.com/jimsyyap/tennis-tracker/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string       `json:"token"`
	User  models.User  `json:"user"`
}

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		RespondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create user
	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	// TODO: Save user to database
	// For now, we'll just respond with success

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return user and token
	RespondWithJSON(w, http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		RespondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// TODO: Get user from database
	// For now, we'll just simulate a user
	user := models.User{
		ID:           1,
		Name:         "Test User",
		Email:        req.Email,
		PasswordHash: "$2a$10$GckdTPS05LnUmMSxH4g4peVKqEoR8h1WUSLfFrS0WkpLcKTVQSfQq", // "password"
	}

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return user and token
	RespondWithJSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

// ForgotPassword handles password reset requests
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate email
	if req.Email == "" {
		RespondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	// TODO: Send password reset email
	// For now, we'll just respond with success

	RespondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "If an account with that email exists, a password reset link has been sent",
	})
}

// ResetPassword handles password reset
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if req.Token == "" || req.Password == "" {
		RespondWithError(w, http.StatusBadRequest, "Token and password are required")
		return
	}

	// TODO: Verify token and update password
	// For now, we'll just respond with success

	RespondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Password has been reset successfully",
	})
}
