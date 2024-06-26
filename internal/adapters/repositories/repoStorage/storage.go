package repoStorage

import (
	"cloud.google.com/go/storage"
	"context"
)

type StorageConfig struct {
	Client *storage.Client
	Bucket *storage.BucketHandle
	Ctx    context.Context
}

type StorageInterface interface {
	Initialize(ctx context.Context, client *storage.Client)
	UploadToStorage(fileName, blob, contentType string) (string, string, error)
}
