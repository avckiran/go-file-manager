package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
)

func main() {
	bucketName := "chandra-cloud-files-data-m512"
	objectPath := "data/input/current/sample_item_1.json"

	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}

	defer client.Close()

	log.Printf("Attempting to read gs://%s/%s\n", bucketName, objectPath)

	content, err := downloadFile(ctx, client, bucketName, objectPath)
	if err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}

	fmt.Println("File Content: ")
	fmt.Println(string(content))

}
