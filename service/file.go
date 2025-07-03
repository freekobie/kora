package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/freekobie/kora/model"
	"github.com/google/uuid"
)

// FileService is a service for managing files.
type FileService struct {
	store model.FileStorage
	gcs   *GCS
}

// NewFileService creates a new FileService.
func NewFileService(store model.FileStorage, gcs *GCS) *FileService {
	return &FileService{store: store, gcs: gcs}
}

type CreateFolderRequest struct {
	Name     string    `json:"name"`
	ParentID uuid.UUID `json:"parentId"`
}

func (s *FileService) CreateFolder(ctx context.Context, req CreateFolderRequest, userID uuid.UUID) (model.Folder, error) {
	folder := model.Folder{
		Id:           uuid.New(),
		Name:         req.Name,
		UserID:       userID,
		ParentID:     req.ParentID,
		CreatedAt:    time.Now(),
		LastModified: time.Now(),
	}

	err := s.store.CreateFolder(ctx, &folder)
	if err != nil {
		return model.Folder{}, err
	}

	return folder, nil
}

// UploadFile uploads a file and saves its metadata.
func (s *FileService) UploadFile(ctx context.Context, userId, folderID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*model.File, error) {
	fileID := uuid.New()
	storageKey := fmt.Sprintf("%s/%s", userId, fileID.String())

	if err := s.gcs.UploadFile(ctx, storageKey, file); err != nil {
		return nil, err
	}

	dbFile := &model.File{
		Id:         fileID,
		Name:       header.Filename,
		UserID:     userId,
		FolderID:   folderID,
		MimeType:   header.Header.Get("Content-Type"),
		Size:       header.Size,
		StorageKey: storageKey,
	}

	if err := s.store.CreateFile(ctx, dbFile); err != nil {
		return nil, err
	}

	return dbFile, nil
}
