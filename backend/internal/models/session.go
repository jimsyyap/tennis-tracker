package models

import (
	"context"
	"time"

	"github.com/jimsyyap/tennis-tracker/backend/internal/database"
)

// Session represents a tennis session
type Session struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Name         string    `json:"name"`
	OpponentName string    `json:"opponent_name,omitempty"`
	SessionDate  time.Time `json:"session_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ErrorCount   int       `json:"error_count,omitempty"` // Total errors for this session
}

// SessionService handles database operations for sessions
type SessionService struct {
	DB *database.DB
}

// GetByID retrieves a session by ID
func (s *SessionService) GetByID(id int) (*Session, error) {
	var session Session
	
	query := `
		SELECT s.id, s.user_id, s.name, s.opponent_name, s.session_date, s.created_at, s.updated_at,
		       COALESCE(SUM(e.count), 0) as error_count
		FROM sessions s
		LEFT JOIN errors e ON s.id = e.session_id
		WHERE s.id = $1
		GROUP BY s.id
	`
	
	err := s.DB.Pool.QueryRow(context.Background(), query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.Name,
		&session.OpponentName,
		&session.SessionDate,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.ErrorCount,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &session, nil
}

// GetByUserID retrieves all sessions for a user
func (s *SessionService) GetByUserID(userID int) ([]Session, error) {
	var sessions []Session
	
	query := `
		SELECT s.id, s.user_id, s.name, s.opponent_name, s.session_date, s.created_at, s.updated_at,
		       COALESCE(SUM(e.count), 0) as error_count
		FROM sessions s
		LEFT JOIN errors e ON s.id = e.session_id
		WHERE s.user_id = $1
		GROUP BY s.id
		ORDER BY s.session_date DESC
	`
	
	rows, err := s.DB.Pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var session Session
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.Name,
			&session.OpponentName,
			&session.SessionDate,
			&session.CreatedAt,
			&session.UpdatedAt,
			&session.ErrorCount,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return sessions, nil
}

// Create inserts a new session into the database
func (s *SessionService) Create(session *Session) error {
	query := `
		INSERT INTO sessions (user_id, name, opponent_name, session_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	
	err := s.DB.Pool.QueryRow(
		context.Background(),
		query,
		session.UserID,
		session.Name,
		session.OpponentName,
		session.SessionDate,
	).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)
	
	return err
}

// Update updates an existing session
func (s *SessionService) Update(session *Session) error {
	query := `
		UPDATE sessions
		SET name = $2, opponent_name = $3, session_date = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`
	
	err := s.DB.Pool.QueryRow(
		context.Background(),
		query,
		session.ID,
		session.Name,
		session.OpponentName,
		session.SessionDate,
	).Scan(&session.UpdatedAt)
	
	return err
}

// Delete removes a session from the database
func (s *SessionService) Delete(id int) error {
	query := `DELETE FROM sessions WHERE id = $1`
	
	_, err := s.DB.Pool.Exec(context.Background(), query, id)
	
	return err
}
