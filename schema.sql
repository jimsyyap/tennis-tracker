-- Tennis Tracker Database Schema

-- Drop tables if they exist (in reverse order of dependencies)
DROP TABLE IF EXISTS shared_links;
DROP TABLE IF EXISTS errors;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create sessions table
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    opponent_name VARCHAR(255),
    session_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create errors table
CREATE TABLE errors (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES sessions(id) ON DELETE CASCADE,
    count INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create shared links table
CREATE TABLE shared_links (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    session_id INTEGER REFERENCES sessions(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_errors_session_id ON errors(session_id);
CREATE INDEX idx_shared_links_token ON shared_links(token);
CREATE INDEX idx_shared_links_user_id ON shared_links(user_id);
CREATE INDEX idx_shared_links_session_id ON shared_links(session_id);

-- Add some sample data (optional, comment out if not needed)
-- Insert sample users
INSERT INTO users (email, password_hash, name) VALUES 
('player1@example.com', '$2a$10$GckdTPS05LnUmMSxH4g4peVKqEoR8h1WUSLfFrS0WkpLcKTVQSfQq', 'Player One'),
('coach1@example.com', '$2a$10$GckdTPS05LnUmMSxH4g4peVKqEoR8h1WUSLfFrS0WkpLcKTVQSfQq', 'Coach One');

-- Insert sample sessions
INSERT INTO sessions (user_id, name, opponent_name, session_date) VALUES 
(1, 'Practice Session 1', 'Practice Partner', NOW() - INTERVAL '7 DAYS'),
(1, 'Tournament Match 1', 'John Doe', NOW() - INTERVAL '3 DAYS');

-- Insert sample errors
INSERT INTO errors (session_id, count) VALUES 
(1, 12),
(2, 8);

-- Insert sample shared link
INSERT INTO shared_links (user_id, session_id, token, expires_at) VALUES 
(1, 1, 'a1b2c3d4e5f6', NOW() + INTERVAL '7 DAYS');

-- Grant permissions (adjust as needed for your environment)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_user;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO your_user;
