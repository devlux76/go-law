// Main processor for orchestrating CSV embedding and storage in Qdrant.
package main

import (
	"log"
	"maps"
	"path/filepath"
	"strings"

	"csv-processor/internal/config"
	"csv-processor/internal/csv"
	"csv-processor/internal/embeddings"
	"csv-processor/internal/qdrant"

	"github.com/spf13/pflag"
)

func main() {
	// CLI flags
	var (
		csvFile          string
		columnEmbed      string
		columnsKept      string
		qdrantHost       string
		qdrantPort       int
		embeddingAPIBase string
		embeddingModel   string
	)

	pflag.StringVar(&csvFile, "csv-file", "", "Path to a single CSV file to process")
	pflag.StringVar(&columnEmbed, "column-embed", "text", "Column name to use for embedding")
	pflag.StringVar(&columnsKept, "columns-kept", "", "Comma-separated list of columns to keep")
	pflag.StringVar(&qdrantHost, "qdrant-host", "localhost", "Qdrant host")
	pflag.IntVar(&qdrantPort, "qdrant-port", 6333, "Qdrant port")
	pflag.StringVar(&embeddingAPIBase, "embedding-api-base", "", "Embedding API base URL (required)")
	pflag.StringVar(&embeddingModel, "embedding-model", "", "Embedding model name (required)")
	pflag.Parse()

	// Load config with flags prioritized
	err := config.LoadWithFlags(qdrantHost, qdrantPort, embeddingAPIBase, embeddingModel)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	var files []string
	if csvFile != "" {
		files = []string{csvFile}
	} else {
		volumePath := "/data"
		files, err = filepath.Glob(filepath.Join(volumePath, "*.csv"))
		if err != nil {
			log.Fatalf("Failed to list CSV files: %v", err)
		}
		if len(files) == 0 {
			log.Printf("No CSV files found in %s", volumePath)
			return
		}
	}

	// Parse columns to keep
	var keptCols []string
	if columnsKept != "" {
		for col := range strings.SplitSeq(columnsKept, ",") {
			c := strings.TrimSpace(col)
			if c != "" {
				keptCols = append(keptCols, c)
			}
		}
	}

	for _, file := range files {
		log.Printf("Processing file: %s", file)
		rows, err := csv.ReadCSVWithDuckDB(file)
		if err != nil {
			log.Printf("Error reading CSV %s: %v", file, err)
			continue
		}
		if len(rows) == 0 {
			log.Printf("No rows found in %s", file)
			continue
		}

		var points []map[string]any
		for i, row := range rows {
			textVal, ok := row[columnEmbed]
			if !ok {
				log.Printf("Row %d missing '%s' column, skipping", i, columnEmbed)
				continue
			}
			textStr, ok := textVal.(string)
			if !ok {
				log.Printf("Row %d '%s' column is not a string, skipping", i, columnEmbed)
				continue
			}
			embedding, err := embeddings.GenerateEmbedding(textStr)
			if err != nil {
				log.Printf("Embedding generation failed for row %d: %v", i, err)
				continue
			}
			point := make(map[string]any)
			// Only keep specified columns if set, else keep all
			if len(keptCols) > 0 {
				for _, k := range keptCols {
					if v, ok := row[k]; ok {
						point[k] = v
					}
				}
			} else {
				maps.Copy(point, row)
			}
			point["embedding"] = embedding
			points = append(points, point)

			if (i+1)%1000 == 0 {
				log.Printf("Processed %d rows in %s", i+1, file)
			}
		}

		if len(points) == 0 {
			log.Printf("No valid rows to store for %s", file)
			continue
		}

		collectionName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		err = qdrant.StoreInQdrant(points, collectionName)
		if err != nil {
			log.Printf("Failed to store data in Qdrant for %s: %v", file, err)
			continue
		}
		log.Printf("Successfully processed and stored %d rows from %s", len(points), file)
	}
}
