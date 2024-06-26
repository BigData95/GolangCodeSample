package dependencies

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"log"
	"main/internal/adapters/repositories/repoStorage"
	"sync"
)

type StorageManager struct {
	client   *storage.Client
	initOnce sync.Once
}

func NewStorageManager() *StorageManager {
	return &StorageManager{}
}

func GetStorageClient(ctx context.Context) (*storage.Client, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create repoStorage client: %w", err)
	}
	return client, nil
}

func (s *StorageManager) Initialize(ctx context.Context, storage repoStorage.StorageInterface) error {
	var err error
	s.initOnce.Do(func() {
		s.client, err = GetStorageClient(ctx)
		if err != nil {
			log.Fatalf("Storage client failed: %v", err)
		}
		storage.Initialize(ctx, s.client)
	})
	return err
}

func (s *StorageManager) Cleanup(ctx context.Context) error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}
