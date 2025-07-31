// Main processor for orchestrating CSV embedding and storage in Qdrant.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"csv-processor/internal/csv"
	"csv-processor/internal/embeddings"
	"csv-processor/internal/qdrant"
)

func main() {
	volumePath := "/data"

	files, err := filepath.Glob(filepath.Join(volumePath, "*.csv"))
	if err != nil {
		log.Fatalf("Failed to list CSV files: %v", err)
	}
	if len(files) == 0 {
		log.Printf("No CSV files found in %s", volumePath)
		return
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

		var points []map[string]interface{}
		for i, row := range rows {
			textVal, ok := row["text"]
			if !ok {
				log.Printf("Row %d missing 'text' column, skipping", i)
				continue
			}
			textStr, ok := textVal.(string)
			if !ok {
				log.Printf("Row %d 'text' column is not a string, skipping", i)
				continue
			}
			embedding, err := embeddings.GenerateEmbedding(textStr)
			if err != nil {
				log.Printf("Embedding generation failed for row %d: %v", i, err)
				continue
			}
			point := make(map[string]interface{})
			for k, v := range row {
				point[k] = v
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