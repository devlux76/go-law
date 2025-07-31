package csv

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/marcboeker/go-duckdb"
)

// RowMap represents a single row as a map of column names to values.
type RowMap map[string]interface{}

// ReadCSVWithDuckDB reads a CSV file using DuckDB's read_csv_auto, with a 100MB row length limit,
// and returns up to 50,000 rows as a slice of RowMap.
func ReadCSVWithDuckDB(filePath string) ([]RowMap, error) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		log.Printf("Error opening DuckDB: %v", err)
		return nil, fmt.Errorf("failed to open DuckDB: %w", err)
	}
	defer db.Close()

	// Set max row size to 100MB (104857600 bytes)
	query := fmt.Sprintf(`
		CREATE OR REPLACE TABLE temp_csv AS
		SELECT * FROM read_csv_auto('%s', 
			MAX_ROW_SIZE=104857600
		);
	`, filePath)

	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Error executing DuckDB CSV read: %v", err)
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	rows, err := db.Query("SELECT * FROM temp_csv LIMIT 50000;")
	if err != nil {
		log.Printf("Error selecting rows from DuckDB: %v", err)
		return nil, fmt.Errorf("failed to select rows: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Printf("Error getting columns from DuckDB: %v", err)
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []RowMap
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			log.Printf("Error scanning row from DuckDB: %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := RowMap{}
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			row[colName] = *val
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Row iteration error in DuckDB: %v", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}