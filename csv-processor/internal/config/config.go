// Package config provides configuration management for the csv-processor service.
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	QdrantHost       string
	QdrantPort       int
	EmbeddingAPIBase string
	EmbeddingModel   string
}

var (
	config   *Config
	loadOnce sync.Once
	loadErr  error
)

// Load reads environment variables, applies defaults, and validates configuration.
// Returns an error if required variables are missing or invalid.
func Load() error {
	loadOnce.Do(func() {
		host := getEnv("QDRANT_HOST", "localhost")
		portStr := getEnv("QDRANT_PORT", "6333")
		apiBase := getEnv("EMBEDDING_API_BASE", "")
		model := getEnv("EMBEDDING_MODEL", "")

		port, err := strconv.Atoi(portStr)
		if err != nil || port <= 0 {
			log.Printf("Invalid QDRANT_PORT: %s", portStr)
			loadErr = fmt.Errorf("invalid QDRANT_PORT: %s", portStr)
			return
		}
		if apiBase == "" {
			log.Printf("EMBEDDING_API_BASE is required")
			loadErr = fmt.Errorf("EMBEDDING_API_BASE is required")
			return
		}
		if model == "" {
			log.Printf("EMBEDDING_MODEL is required")
			loadErr = fmt.Errorf("EMBEDDING_MODEL is required")
			return
		}

		config = &Config{
			QdrantHost:       host,
			QdrantPort:       port,
			EmbeddingAPIBase: apiBase,
			EmbeddingModel:   model,
		}
	})
	return loadErr
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

// GetQdrantHost returns the configured Qdrant host.
func GetQdrantHost() string {
	return config.QdrantHost
}

// GetQdrantPort returns the configured Qdrant port.
func GetQdrantPort() int {
	return config.QdrantPort
}

// GetEmbeddingAPIBase returns the configured embedding API base URL.
func GetEmbeddingAPIBase() string {
	return config.EmbeddingAPIBase
}

// GetEmbeddingModel returns the configured embedding model name.
func GetEmbeddingModel() string {
	return config.EmbeddingModel
}