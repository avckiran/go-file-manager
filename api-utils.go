package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func fetchAPIData(ctx context.Context, url string, target interface{}) error {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return fmt.Errorf("failed to create API request, %w", err)
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("failed to execute API request, %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %s: %s", resp.Status, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to read API response body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal API response JSON: %w", err)
	}

	return nil

}
