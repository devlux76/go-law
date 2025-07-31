# Docker Container Design for CSV Processing

## Overview
This design outlines a Docker container setup for processing CSV files from the host's `/srv/bulk-data` directory. The container will use DuckDB to read CSV files, generate embeddings for text fields using an OpenAI-compatible API, and store the results in Qdrant.

We will be developing and testing with the "/srv/bulk-data/parentheticals-latest.csv" 
It has the following header row
```
id,text,score,described_opinion_id,describing_opinion_id,group_id
```

## Components

### 1. Docker Container
- **Base Image**: Go latest (golang:tip-bookworm) with necessary libraries and bindings (qdrant and duckdb) 
- **Mounted Volume**: `/srv/bulk-data` from host
- **Environment Variables**:
  - `QDRANT_HOST`: Hostname for Qdrant service
  - `QDRANT_PORT`: Port for Qdrant service
  - `EMBEDDING_API_BASE`: Base URL for embedding API
  - `EMBEDDING_MODEL`: Model name for embeddings

### 2. Processing Workflow
1. **CSV Reading**: Use DuckDB's `read_csv_auto` with 100MB row length
2. **Data Selection**: Select all columns (default 50k rows)
3. **Embedding Generation**:
   - Use "nomic-embed-text" model
   - API endpoint on the host already running at: `localhost:2002/v1`
   - Format text as "search_document: <text here>" for high-quality embeddings
   - The embedding model requires text to be formatted as "search_document: <text here>"
   - This prefix ensures high-quality embedding generation
   - The text from the CSV's "text" column should be prepended with "search_document: " before sending to the embedding API
   - Generate embeddings for the "text" column with this prefix
4. **Data Storage**:
   - Convert rows to JSON
   - Store in Qdrant using filename (without extension) as collection name
   - Example: "parentheticals-latest.csv" → "parentheticals-latest" collection
   - qdrant is already running on the host at localhost:6333

### 3. Docker Compose
- **Services**:
  - `csv-processor`: Main processing container
  qdrant and nomic-embed-text are already running on the host in docker form but may need to be added to the network.
- **Networking**: Ensure all services can communicate

### 4. Extensibility
- Design should be general enough to handle multiple CSV files
- Configuration should allow easy addition of new CSV files
- Processing logic should be modular for different text columns

## Implementation Notes
- The embedding API must be accessible from within the container, it will be an openai compatible API
- Qdrant service should be properly configured for vector storage
- Error handling for API failures and CSV parsing issues
- Logging for monitoring progress and debugging