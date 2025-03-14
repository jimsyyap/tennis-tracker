package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never send to client
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserService handles database operations for users
type UserService struct {
	DB *database.DB
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id int) (*User, error) {
	var user User
	
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	
	err := s.DB.Pool.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(email string) (*User, error) {
	var user User
	
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	
	err := s.DB.Pool.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// Create inserts a new user into the database
func (s *UserService) Create(user *User) error {
	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	
	err := s.DB.Pool.QueryRow(
		context.Background(),
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	
	return err
}

// Update updates an existing user
func (s *UserService) Update(user *User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`
	
	err := s.DB.Pool.QueryRow(
		context.Background(),
		query,
		user.ID,
		user.Name,
		user.Email,
	).Scan(&user.UpdatedAt)
	
	return err
}

// UpdatePassword updates a user's password
func (s *UserService) UpdatePassword(userID int, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $2, updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := s.DB.Pool.Exec(
		context.Background(),
		query,
		userID,
		passwordHash,
	)
	
	return err
}
