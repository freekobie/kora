package postgres

import (
	"context"
	"log/slog"

	"github.com/freekobie/kora/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FileStore is a repository for managing files in a PostgreSQL database.
type FileStore struct {
	conn *pgxpool.Pool
}

// NewFileStore creates a new FileRepository.
func NewFileStore(conn *pgxpool.Pool) model.FileStorage {
	return &FileStore{conn: conn}
}

func (s *FileStore) CreateFolder(ctx context.Context, folder *model.Folder) error {
	query := `
		INSERT INTO folders (id, name, user_id, parent_id, created_at, last_modified)
		VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := s.conn.Exec(ctx, query,
		folder.Id,
		folder.Name,
		folder.UserID,
		folder.ParentID,
		folder.CreatedAt,
		folder.LastModified,
	)

	if err != nil {
		slog.Error("failed to insert folder", "error", err)
		return err
	}

	return nil
}

// CreateFile creates a new file in the database.
func (r *FileStore) CreateFile(ctx context.Context, file *model.File) error {
	query := `
		INSERT INTO files (id, name, user_id, folder_id, mime_type, size, storage_key)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.conn.Exec(ctx, query, file.Id, file.Name, file.UserID, file.FolderID, file.MimeType, file.Size, file.StorageKey)
	return err
}
