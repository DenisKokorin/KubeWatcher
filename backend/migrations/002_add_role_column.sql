-- +goose Up
-- +goose StatementBegin

-- Add role column to users table
ALTER TABLE users ADD COLUMN role VARCHAR(50) NOT NULL DEFAULT 'user';

-- Create index on role for faster queries
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_users_role;
ALTER TABLE users DROP COLUMN role;

-- +goose StatementEnd
