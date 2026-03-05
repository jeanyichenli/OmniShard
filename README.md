<!-- @format -->

# OmniShard

## Introduction

### Distributed File Upload & Integrity Verification System

A modular backend system for handling chunked file uploads with resumable support, cryptographic verification (SHA256 / AES-GCM), and distributed metadata tracking.

In the recent refactoring phase, I adopted Conventional Commits to improve project maintainability.

### Project structure (current)

```
OmniShard/
├── cmd/           # CLI entrypoints (serve, upload, version)
├── api/           # HTTP server, Gin routes and handlers
├── internal/      # Internal domain services (e.g. upload service)
├── chunker/       # Chunk-size calculation and chunking helpers
├── datastore/     # File metadata / storage abstraction
├── redis/         # Redis client initialization
├── logger/        # Logging utilities
├── config         # Configuration stub / helpers
├── main.go        # Main entrypoint (wires CLI + Redis)
├── go.mod, go.sum
└── README.md
```

### Features

- Chunk-based file upload
- Data storage in disk
- Metadata storage in Redis / MongoDB
- Crypto method + SHA256 for content integrity
- REST API for upload status and control
- CLI upload tool
- Containerized & K8s practices

---

## Quick start

### Environment requirements

- **Go**: 1.23 or later (see `go.mod`)
- **Redis**: reachable from the application, configured via `REDIS_URI`
  - Example: `redis://localhost:6379`

### Configuration

- **Redis connection**
  - Set `REDIS_URI` before running the binary so `redis.InitRedisClient()` can connect:
    - `export REDIS_URI="redis://localhost:6379"`

### Install

From the project root:

```bash
# Build the OmniShard binary into ./bin/OmniShard
make build

# (Optional) Install to /usr/local/bin/OmniShard
sudo make install
```

After installation you can run the CLI as `OmniShard`.

### Usage

#### 1. Start the HTTP server

Start the Gin HTTP server (default port 8080):

```bash
OmniShard serve
```

This will expose:

- `POST /upload` – accepts an upload session JSON payload.

#### 2. Trigger an upload from the CLI

Use the `upload` command to send an upload session to the HTTP API:

```bash
OmniShard upload \
  --file-path /path/to/file.bin \
  --chunk-size 256k
```

The CLI will:

- Open the file and calculate:
  - total size
  - chunk size (based on `--chunk-size`, e.g. `256k`, `1M`)
  - total number of chunks
- Build an `UploadSession` JSON body.
- `POST` it to `http://localhost:8080/upload` with `Content-Type: application/json`.

On success, the API responds with a JSON payload similar to:

```json
{
  "uploadid": "<uuid>",
  "status": "INITIATED"
}
```

This `uploadid` can be used later for tracking the upload session as more features are added.

## Future work

Planned next steps (not implemented yet):

- **HTTP/API layer**
  - Endpoints to check the status of upload processing with uploadid
  - Endpoints to download the uploaded file
- **Additional features**
  - Logging service for recording code events
  - Metadata structure stored in MongoDB
  - Dockerfile and Docker Compose setup for containerized and K8s deployment
