package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jimsyyap/tennis-tracker/backend/internal/database"
	customMiddleware "github.com/jimsyyap/tennis-tracker/backend/internal/middleware"
)

// NewRouter sets up and returns the router for the API
func NewRouter(db *database.DB) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	
	// Custom middleware
	r.Use(customMiddleware.JSONContentType)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/health", HealthCheck)
		
		// Auth endpoints
		r.Post("/api/register", Register)
		r.Post("/api/login", Login)
		r.Post("/api/forgot-password", ForgotPassword)
		r.Post("/api/reset-password", ResetPassword)
		
		// Shared data endpoint (public)
		r.Get("/api/shared/{token}", GetSharedSession)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		// Use authentication middleware
		r.Use(customMiddleware.Authenticate)
		
		// User endpoints
		r.Get("/api/user", GetUser)
		r.Put("/api/user", UpdateUser)
		
		// Session endpoints
		r.Route("/api/sessions", func(r chi.Router) {
			r.Get("/", GetSessions)
			r.Post("/", CreateSession)
			r.Get("/{id}", GetSession)
			r.Put("/{id}", UpdateSession)
			r.Delete("/{id}", DeleteSession)
			
			// Error tracking endpoints
			r.Route("/{sessionID}/errors", func(r chi.Router) {
				r.Get("/", GetErrors)
				r.Post("/", LogError)
				r.Put("/{id}", UpdateError)
				r.Delete("/{id}", DeleteError)
			})
			
			// Sharing endpoints
			r.Post("/{id}/share", ShareSession)
			r.Delete("/{id}/share", RemoveShare)
		})
	})

	return r
}

// HealthCheck handles the health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
