package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// GCS is a service for interacting with Google Cloud Storage.
type GCS struct {
	client *storage.Client
	bucket string
}

// NewGCS creates a new GCS service.
func NewGCS(ctx context.Context, bucket string) (*GCS, error) {
	// client, err := storage.NewClient(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create gcs client: %w", err)
	// }

	// return &GCS{
	// 	client: client,
	// 	bucket: bucket,
	// }, nil

	return &GCS{}, nil
}

// UploadFile uploads a file to Google Cloud Storage.
func (s *GCS) UploadFile(ctx context.Context, key string, file io.Reader) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := s.client.Bucket(s.bucket).Object(key).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("failed to copy file to gcs: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close gcs writer: %w", err)
	}

	return nil
}
