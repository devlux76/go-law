# CSV Processor

A robust and extensible tool for processing large CSV datasets, designed for bulk data workflows. This project leverages Go for high performance and is containerized for easy deployment and reproducibility.

## Prerequisites

- [Docker](https://www.docker.com/get-started) installed on your system
- Access to the CSV files you wish to process (recommended location: `/srv/bulk-data` on your host)

## Installation

Clone this repository and navigate to the `csv-processor` directory:

```bash
git clone https://github.com/your-org/go-law.git
cd go-law/csv-processor
```

## Usage with Docker

The recommended way to run the CSV Processor is via Docker. This ensures all dependencies are handled and the environment is consistent.

### 1. Build the Docker Image

```bash
docker build -t csv-processor .
```

### 2. Run the Container

Mount your bulk data directory (`/srv/bulk-data`) into the container at the same path:

```bash
docker run --rm -v /srv/bulk-data:/srv/bulk-data csv-processor
```

- The container expects CSV files to be available in `/srv/bulk-data`.
- Processed output and logs will also be written to `/srv/bulk-data` unless otherwise configured.

### 3. Using Docker Compose (Optional)

A `docker-compose.yml` is provided for convenience. To use it:

```bash
docker-compose up --build
```

This will build and start the service with the appropriate volume mount.

## Mounting `/srv/bulk-data`

- **Required:** The `/srv/bulk-data` directory must be mounted into the container for the processor to access input files and write outputs.
- **Example:**  
  `-v /srv/bulk-data:/srv/bulk-data`
- Ensure your CSV files are placed in `/srv/bulk-data` on the host before running the container.

## Configuration

- Configuration options can be set via environment variables or by editing the configuration files in the project.
- See [`internal/config/config.go`](csv-processor/internal/config/config.go:1) for available configuration parameters.

## Project Structure

- [`cmd/processor/main.go`](csv-processor/cmd/processor/main.go:1): Entry point for the processor.
- [`internal/csv/reader.go`](csv-processor/internal/csv/reader.go:1): CSV reading logic.
- [`internal/embeddings/generator.go`](csv-processor/internal/embeddings/generator.go:1): Embedding generation.
- [`internal/qdrant/client.go`](csv-processor/internal/qdrant/client.go:1): Qdrant vector database client.

## Contributing

Contributions are welcome! Please open issues or submit pull requests for improvements or bug fixes.

## License

This project is licensed under the MIT License.