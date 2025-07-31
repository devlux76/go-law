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
	FromFlags        bool
}

var (
	config   *Config
	loadOnce sync.Once
	loadErr  error
)

/*
LoadWithFlags loads config, prioritizing CLI flags over environment variables.
If a flag value is set (not default), it is used; otherwise, env var is used.
*/
func LoadWithFlags(flagHost string, flagPort int, flagAPIBase string, flagModel string) error {
	loadOnce.Do(func() {
		host := flagHost
		if host == "" || host == "localhost" {
			host = getEnv("QDRANT_HOST", "localhost")
		}
		port := flagPort
		if port == 0 || port == 6333 {
			portStr := getEnv("QDRANT_PORT", "6333")
			p, err := strconv.Atoi(portStr)
			if err == nil && p > 0 {
				port = p
			}
		}
		apiBase := flagAPIBase
		if apiBase == "" {
			apiBase = getEnv("EMBEDDING_API_BASE", "")
		}
		model := flagModel
		if model == "" {
			model = getEnv("EMBEDDING_MODEL", "")
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
			FromFlags:        true,
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