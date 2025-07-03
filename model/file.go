package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Folder struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	UserID       uuid.UUID `json:"userId"`
	ParentID     uuid.UUID `json:"parentId,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	LastModified time.Time `json:"lastModified"`
}

type File struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	UserID       uuid.UUID `json:"userId"`
	FolderID     uuid.UUID `json:"folderId,omitempty"`
	MimeType     string    `json:"mimeType"`
	Size         int64     `json:"size"`
	StorageKey   string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	LastModified time.Time `json:"lastModified"`
}


// FileStorage is an interface for storing and retrieving file metadata.
type FileStorage interface {
	CreateFolder(ctx context.Context, folder *Folder) error
	CreateFile(ctx context.Context, file *File) error
}

