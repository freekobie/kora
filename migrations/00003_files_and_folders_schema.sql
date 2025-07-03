-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS folders (
    id uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id uuid REFERENCES folders(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    last_modified TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS files (
    id uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    folder_id uuid REFERENCES folders(id) ON DELETE CASCADE,
    mime_type VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    storage_key TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    last_modified TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_folders_user_id ON folders (user_id);
CREATE INDEX idx_files_user_id ON files (user_id);
CREATE INDEX idx_files_folder_id ON files (folder_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS folders;
-- +goose StatementEnd
