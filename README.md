# Project catalog-be

Monolith REST-API service for `Inner Catalog` project.

# Community

- Should i make a discord server?

## Getting Started

## Docs

- Use [Bruno](https://www.usebruno.com/)
  - Import collection form `./docs/catalog-circle-api/bruno.json`

### Requirements

- Go 1.21.1
- PostgreSQL:latest
- air (for live reload)
- Cloudflare R2 / Object Storage
- Google Client ID & Secret (for oauth2)
- Docker

### Installation

1. Fork the repository
2. Clone the repository from your fork
3. Create `.env` file in the root dir based on `.env.example`
4. Install dependencies

```bash
go mod download
```

5. Run DB container

```bash
make docker-run
```

6. Run the application locally

```bash
make watch
```

## Available Make Commands

```bash
# run all make commands with clean tests
make all build

# build the application
make build

# run the application
make run

# Create DB container
make docker-run

# Shutdown DB container
make docker-down

# live reload the application
make watch

# run the test suite
make test

# clean up binary from the last build
make clean
```
