package repoStorage

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

type TrackingStorageImpl struct {
	config *StorageConfig
}

func (t *TrackingStorageImpl) Initialize(ctx context.Context, client *storage.Client) {
	bucket := client.Bucket(os.Getenv("DEFAULT_BUCKET_NAME"))
	t.config = &StorageConfig{
		Client: client,
		Bucket: bucket,
		Ctx:    ctx,
	}
}

func (t *TrackingStorageImpl) UploadToStorage(fileName, blob, contentType string) (shipperReference, fileExtension string, err error) {
	o := t.config.Client.Bucket(os.Getenv("DEFAULT_BUCKET_NAME")).Object(blob)
	o = o.If(storage.Conditions{DoesNotExist: true})

	fileExtension = strings.Split(contentType, "/")[1]
	path := fmt.Sprintf("evidences/%s.%s", fileName, fileExtension)

	wc := t.config.Bucket.Object(path).NewWriter(t.config.Ctx)
	wc.ContentType = contentType

	rawDecodedBlob, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode blob: %w", err)
	}

	if _, err := wc.Write(rawDecodedBlob); err != nil {
		return "", "", fmt.Errorf("failed to write to bucket: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", "", fmt.Errorf("failed to close writer: %w", err)
	}

	return fmt.Sprintf("gs://%s/%s", os.Getenv("DEFAULT_BUCKET_NAME"), path), fileExtension, nil

}
