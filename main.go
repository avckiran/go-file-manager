package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

func archiveMasterFileIfExists(ctx context.Context, client *storage.Client, bucketName, masterOutputPath string) {
	exists, err := objectExists(ctx, client, bucketName, masterOutputPath)

	if err != nil {
		log.Fatalf("Failed to check for master file existence: %v", err)
	}

	if exists {
		log.Printf("Existing master file found at %s. Archiving...", masterOutputPath)

		archiveDate := time.Now().Format("20060102_150405")
		archiveFileName := fmt.Sprintf("final_file_%s.xlsx", archiveDate)
		archivePath := filepath.Join("data", "output", "archive", archiveFileName)
		archivePath = strings.ReplaceAll(archivePath, "\\", "/")

		log.Printf("ARchive path %s", archivePath)

		err := copyGCSObject(ctx, client, bucketName, masterOutputPath, archivePath)

		if err != nil {
			log.Fatalf("Failed to copy master file to archive: %v", err)
		}

		err = deleteGCSObject(ctx, client, bucketName, masterOutputPath)
		if err != nil {
			log.Fatalf("Failed to delete master file from the original path: %v", err)
		}

		log.Printf("Successfully archived the existing master file")
	} else {
		log.Println("No existing masterfile found")
	}
}

func main() {

	// Retrieving json file from GCS bucket

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

	var inputRecord InputData

	if err := json.Unmarshal(content, &inputRecord); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	log.Println("Successfully unmarshalled JSON data!")

	// API Data fetch for Cat Facts

	apiUrl := "https://catfact.ninja/fact"

	var catFactData ApiResponseData

	if err := fetchAPIData(ctx, apiUrl, &catFactData); err != nil {
		log.Fatalf("Failed to retrieve api data %v", err)
	}
	log.Println("Successfully fetched and unmarshalled API data!")
	fmt.Printf("Cat Fact: %s (Length: %d)\n", catFactData.Fact, catFactData.Length)

	log.Println("Aggregating the data")

	aggregatedData := AggregatedRecord{
		ItemID:         inputRecord.ItemId,
		ItemName:       inputRecord.ItemName,
		Quantity:       inputRecord.Quantity,
		Color:          inputRecord.Attributes.Color,
		Size:           inputRecord.Attributes.Size,
		FetchedAPIData: catFactData.Fact,
		APISource:      "Cat Facts (Cat Ninja)",
		ProcessingDate: time.Now().Format("2006-01-02"),
	}

	log.Println("Aggregation complete")
	fmt.Printf("Aggregated Data: %+v\n", aggregatedData)

	// Moving processed file to another location in GCS

	log.Println("Moving processed file in GCS...")

	_, destinationObjectPath := getProcessedFileName(objectPath)

	err = copyGCSObject(ctx, client, bucketName, objectPath, destinationObjectPath)

	if err != nil {
		log.Printf("Error copying object, proceeding without deleting original: %v", err)
	} else {
		err = deleteGCSObject(ctx, client, bucketName, objectPath)
		if err != nil {
			log.Printf("Unable to delete the object %v", err)
		} else {
			log.Printf("Successfully moved the file to the processed folder!")
		}
	}

	log.Println("Checking if a master excel file already exists!")

	masterOutputPath := "data/output/master/final_file.xlsx"

	archiveMasterFileIfExists(ctx, client, bucketName, masterOutputPath)

	log.Println("Creating aggregated excel file")

	excelBytes, err := createExcelFile(aggregatedData)

	if err != nil {
		log.Fatalf("Failed to create excel file bytes: %v", err)
	}
	log.Printf("Excel file created in memory %d bytes", len(excelBytes))

	err = uploadFile(ctx, client, bucketName, masterOutputPath, excelBytes)

	if err != nil {
		log.Fatalf("Failed to upload excel file : %v", err)
	}

	log.Println("-- Pipeline finished successfully --")

}
