# Transfer Service

A service for transferring files between nodes.

## Configuration

The service can be configured using environment variables:

- `PORT`: Server port (default: 8080)
- `BASE_URL`: Base URL for the service (default: http://localhost:8080)
- `LOG_LEVEL`: Logging level (default: info)

Example:
```bash
export PORT=9090
export LOG_LEVEL=debug
```

## Development

This project uses Task for managing development tasks. Available tasks:

- `task build`: Build the application
- `task test`: Run tests
- `task lint`: Run linter (golangci-lint)
- `task format`: Format Go code
- `task docker`: Build Docker image
- `task run`: Run the application
- `task dev`: Run in development mode
- `task clean`: Clean build artifacts

### Running the Service

Build and run:
```bash
# Build and run directly
go build -o transfer-service ./cmd/server
./transfer-service

# Or using Task
task run
```

### Development Mode

For development:
```bash
task dev
```

### Testing

Run tests:
```bash
task test
```

### Linting

Run linter:
```bash
task lint
```

## API Endpoints

- `GET /health`: Health check endpoint
- `POST /jobs`: Create a new transfer job
- `GET /jobs/{id}`: Get job status
- `DELETE /jobs/{id}`: Cancel a job

## Docker

Build and run using Docker:
```bash
task docker
docker run -p 8080:8080 transfer-service
``` 