# Project catalog-api 📚

Monolith REST API service for `Inner Catalog` project.

## Community

- Just use github discusson for now
- Should i make a discord server?

## Docs

- Use [Bruno](https://www.usebruno.com/)
  - Import collection form `./docs/catalog-circle-api/bruno.json`
- [ERD Diagram](./docs/erd/README.md)

## Getting Started

### Requirements

- Go 1.21.1
- PostgreSQL:latest
  - Docker (local)
  - Supabase
- Go Migrate v4 [link](https://github.com/golang-migrate/migrate)
- wire
- Docker
- air (for live reload)
- Google Client ID & Secret (for oauth2)
- Object Storage
  - AWS S3 with Cloudflare R2

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

6. Run the migration

```bash
make migrate-up
```

6. Run the application locally

```bash
make watch
```

## Environment

- dev - development environment [https://api-dev.innercatalog.com](https://api-dev.innercatalog.com)
- prod - production environment [https://api.innercatalog.com](https://api.innercatalog.com) (SOON)

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

# create new migration file
make migrate-create

# run the migration
make migrate-up

# rollback the migration
make migrate-down

# run the test suite
make test

# clean up binary from the last build
make clean
```
