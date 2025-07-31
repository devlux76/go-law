// Package qdrant provides a client for storing data in Qdrant.
package qdrant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// StoreInQdrant stores the given data in Qdrant under the specified collection name.
// data should be a slice of maps or structs representing the points to store.
// collectionName is the Qdrant collection (use filename without extension).
func StoreInQdrant(data interface{}, collectionName string) error {
	// Marshal data to JSON
	payload := map[string]interface{}{
		"points": data,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal data to JSON for Qdrant: %v", err)
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	url := fmt.Sprintf("http://localhost:6333/collections/%s/points", collectionName)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to create HTTP request for Qdrant: %v", err)
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request to Qdrant: %v", err)
		return fmt.Errorf("failed to send request to Qdrant: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Qdrant returned non-2xx status: %d", resp.StatusCode)
		return fmt.Errorf("Qdrant returned status %d", resp.StatusCode)
	}

	return nil
}