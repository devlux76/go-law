// Package embeddings provides embedding generation functionality.
package embeddings

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

// EmbeddingAPIURL is the endpoint for the embedding model.
const EmbeddingAPIURL = "http://localhost:2002/v1/embeddings"

// EmbeddingRequest is the request payload for the embedding API.
type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// EmbeddingResponse is the response from the embedding API.
type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// GenerateEmbedding generates an embedding for the given text using the nomic-embed-text model.
func GenerateEmbedding(text string) ([]float32, error) {
	formatted := "search_document: " + text

	reqBody := EmbeddingRequest{
		Model: "nomic-embed-text",
		Input: []string{formatted},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Error marshaling embedding request: %v", err)
		return nil, fmt.Errorf("failed to marshal embedding request: %w", err)
	}

	resp, err := http.Post(EmbeddingAPIURL, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		log.Printf("Embedding API request failed: %v", err)
		return nil, fmt.Errorf("embedding API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read embedding API response: %v", err)
		return nil, fmt.Errorf("failed to read embedding API response: %w", err)
	}

	var embeddingResp EmbeddingResponse
	if err := json.Unmarshal(respBody, &embeddingResp); err != nil {
		log.Printf("Failed to unmarshal embedding API response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal embedding API response: %w", err)
	}

	if embeddingResp.Error != nil {
		log.Printf("Embedding API returned error: %s", embeddingResp.Error.Message)
		return nil, errors.New(embeddingResp.Error.Message)
	}

	if len(embeddingResp.Data) == 0 {
		log.Printf("Embedding API returned no data for input: %s", text)
		return nil, errors.New("embedding API returned no data")
	}

	return embeddingResp.Data[0].Embedding, nil
}