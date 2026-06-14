# Bookstore API

A simple Go REST API backed by SQLite for books, authors, customers, and purchases.

## Features

- CRUD endpoints for `books`, `authors`, and `customers`
- Purchase tracking via `customer_book_purchase`
- SQLite database with migrations
- Unit tests and integration tests
- GitHub Actions CI with formatting, testing, and code review

## Endpoints

- `GET /books`
- `GET /books/{id}`
- `POST /books`
- `PUT /books/{id}`
- `DELETE /books/{id}`

- `GET /authors`
- `GET /authors/{id}`
- `POST /authors`
- `PUT /authors/{id}`
- `DELETE /authors/{id}`

- `GET /customers`
- `GET /customers/{id}`
- `POST /customers`
- `PUT /customers/{id}`
- `DELETE /customers/{id}`

- `GET /purchases`
- `POST /purchases`

## Getting Started

### Requirements

- Go 1.22+
- SQLite

### Run locally

```bash
cd /workspaces/codespaces-blank/bookstore-api
go run ./cmd/server
```

The server listens on `:8080` by default.

### Environment

- `PORT` to set the HTTP port
- `DATABASE_PATH` to set the SQLite file path

## Testing

Run unit tests:

```bash
./scripts/run-tests.sh
```

Run integration tests:

```bash
./scripts/run-integration-tests.sh
```

Run all tests:

```bash
./scripts/run-all-tests.sh
```

Run code review / static analysis:

```bash
./scripts/run-code-review.sh
```

## CI

This repository includes a GitHub Actions workflow in `.github/workflows/ci.yml` that:

- checks out the repository
- sets up Go
- formats code with `gofmt`
- runs `go test ./...`
- executes the code review script
- builds the server binary

## Notes

- The `scripts/run-code-review.sh` script runs the reusable Copilot prompt and `golangci-lint`.
- Database migrations are executed automatically on server startup.
