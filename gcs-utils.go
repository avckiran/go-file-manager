package main

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

func downloadFile(ctx context.Context, client *storage.Client, bucketName, objectName string) ([]byte, error) {
	obj := client.Bucket(bucketName).Object(objectName)

	rc, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("obj.NewReader: %w", err)
	}

	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	return data, nil
}
