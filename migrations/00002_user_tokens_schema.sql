-- +goose Up
-- +goose StatementBegin
DO
$$
    BEGIN
        CREATE TYPE scope_type AS ENUM ('authentication', 'verification');
    EXCEPTION
        WHEN duplicate_object THEN null;
    END
$$;

CREATE TABLE IF NOT EXISTS user_tokens (
    token_hash bytea NOT NULL,
    user_id uuid NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    scope scope_type NOT NULL,
    PRIMARY KEY (token_hash),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP type IF EXISTS scope_type;
DROP TABLE IF EXISTS user_tokens;
-- +goose StatementEnd
