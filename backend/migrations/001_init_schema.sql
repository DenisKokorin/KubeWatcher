-- +goose Up
-- +goose StatementBegin

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    uid UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    position VARCHAR(255) NOT NULL,
    team VARCHAR(255) NOT NULL,
    hire_date TIMESTAMP NOT NULL,
    hashed_password VARCHAR(255) NOT NULL
);

-- Create reviews table
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL,
    reviewer_id UUID NOT NULL,
    period VARCHAR(50) NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comments TEXT,
    goals TEXT[],
    strengths TEXT[],
    areas_for_improvement TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (employee_id) REFERENCES users(uid),
    FOREIGN KEY (reviewer_id) REFERENCES users(uid)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_team ON users(team);
CREATE INDEX IF NOT EXISTS idx_reviews_employee_id ON reviews(employee_id);
CREATE INDEX IF NOT EXISTS idx_reviews_reviewer_id ON reviews(reviewer_id);
CREATE INDEX IF NOT EXISTS idx_reviews_period ON reviews(period);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_reviews_period;
DROP INDEX IF EXISTS idx_reviews_reviewer_id;
DROP INDEX IF EXISTS idx_reviews_employee_id;
DROP INDEX IF EXISTS idx_users_team;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
