package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

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

func copyGCSObject(ctx context.Context, client *storage.Client, bucketName, srcObjectName, desObjectName string) error {

	src := client.Bucket(bucketName).Object(srcObjectName)
	des := client.Bucket(bucketName).Object(desObjectName)

	if _, err := des.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("failed to copy object from %s to %s: %w", srcObjectName, desObjectName, err)
	}
	log.Printf("Successfully copied gs://%s/%s to gs://%s/%s\n", bucketName, srcObjectName, bucketName, desObjectName)

	return nil
}

func deleteGCSObject(ctx context.Context, client *storage.Client, bucketName, objectName string) error {

	obj := client.Bucket(bucketName).Object(objectName)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object %s: %w", objectName, err)
	}
	log.Printf("Successfully deleted gs://%s/%s\n", bucketName, objectName)

	return nil
}

func getProcessedFileName(originalPath string) (string, string) {
	fileNameWithExt := filepath.Base(originalPath)
	extension := filepath.Ext(fileNameWithExt)
	fileNameOnly := strings.TrimSuffix(fileNameWithExt, extension)

	dateSuffix := time.Now().Format("20060102")

	processedFileName := fmt.Sprintf("%s_%s%s", fileNameOnly, dateSuffix, extension)
	destinationPath := filepath.Join("data", "input", "processed", processedFileName)

	return processedFileName, strings.ReplaceAll(destinationPath, "\\", "/")

}

func objectExists(ctx context.Context, client *storage.Client, bucketName, objectName string) (bool, error) {
	obj := client.Bucket(bucketName).Object(objectName)

	_, err := obj.Attrs(ctx)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, storage.ErrObjectNotExist) {
		return false, nil
	}

	return false, fmt.Errorf("failed to check for object %s: %w", objectName, err)
}

func uploadFile(ctx context.Context, client *storage.Client, bucketName, objectPath string, data []byte) error {
	obj := client.Bucket(bucketName).Object(objectPath)
	wc := obj.NewWriter(ctx)

	if _, err := wc.Write(data); err != nil {
		return fmt.Errorf("failed to write data to GCS object %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close GCS Object writer: %w", err)
	}

	log.Printf("Successfully uploaded file to gs: //%s/%s\n", bucketName, objectPath)

	return nil
}
