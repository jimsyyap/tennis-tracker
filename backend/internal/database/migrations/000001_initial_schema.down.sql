-- Tennis Tracker Initial Schema Rollback

-- Drop indexes
DROP INDEX IF EXISTS idx_shared_links_session_id;
DROP INDEX IF EXISTS idx_shared_links_user_id;
DROP INDEX IF EXISTS idx_shared_links_token;
DROP INDEX IF EXISTS idx_errors_session_id;
DROP INDEX IF EXISTS idx_sessions_user_id;

-- Drop tables (in reverse order of creation to handle dependencies)
DROP TABLE IF EXISTS shared_links;
DROP TABLE IF EXISTS errors;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
