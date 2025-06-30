-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash bytea NOT NULL,
    profile_photo TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now (),
    last_modified TIMESTAMP NOT NULL DEFAULT now (),
    verified BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_users_email ON users (email);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

-- +goose StatementEnd