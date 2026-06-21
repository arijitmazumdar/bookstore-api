# Repository Guidelines

## Project Structure & Module Organization

This is a Go 1.22 REST API for a SQLite-backed bookstore. The executable entry point is `cmd/server/main.go`. Application wiring, routing, and integration tests live in `internal/app/`. Database connection and migrations are in `internal/db/`. HTTP handlers are grouped by resource in `internal/handlers/`, and shared domain structs are in `internal/models/`. Helper scripts are in `scripts/`, and CI plus reusable agent prompts are under `.github/`.

## Build, Test, and Development Commands

- `go run ./cmd/server`: starts the API on `:8080` by default.
- `PORT=9090 DATABASE_PATH=/tmp/bookstore.db go run ./cmd/server`: runs locally with explicit port and database path.
- `go test ./...`: runs all Go tests.
- `./scripts/run-tests.sh`: wrapper for the full Go test suite.
- `./scripts/run-integration-tests.sh`: runs `TestBookstoreIntegration` in `internal/app`.
- `./scripts/run-all-tests.sh`: runs unit and integration test scripts.
- `go build -o bookstore-api ./cmd/server`: builds the server binary.
- `./scripts/run-code-review.sh`: runs the Copilot review prompt and `golangci-lint`; it may install `golangci-lint` if missing.

## Coding Style & Naming Conventions

Use standard Go formatting with `gofmt`; CI currently runs `gofmt -w .`. Keep packages small and resource-oriented: handlers should remain in `internal/handlers`, database concerns in `internal/db`, and routing in `internal/app`. Prefer clear Go names such as `NewRouter`, `RunMigrations`, and `CreateBookHandler`. Test helpers should call `t.Helper()` and keep temporary SQLite files isolated.

## Testing Guidelines

Tests use Go's built-in `testing` package and `httptest` for HTTP integration coverage. Place tests beside the package under test with `_test.go` suffixes, such as `handlers_test.go` or `integration_test.go`. Name tests after behavior, for example `TestBookstoreIntegration`. Run `go test ./...` before opening a PR, and add or update tests when changing routes, JSON payloads, migrations, or SQL behavior.

## Commit & Pull Request Guidelines

Recent commits use short, imperative, lower-case summaries such as `add README.md with project overview, features, and usage instructions`. Keep new commit subjects direct and scoped. Pull requests should describe the API or behavior change, list tests run, link related issues when applicable, and call out migration or configuration changes. Include example requests or responses when changing endpoints.

## Security & Configuration Tips

Do not commit generated SQLite database files, secrets, or local environment overrides. Configure runtime behavior with `PORT` and `DATABASE_PATH`. Treat SQL changes carefully: keep migrations deterministic, preserve referential integrity, and verify them with integration tests.

<!-- SPECKIT START -->
For additional context about technologies to be used, project structure,
shell commands, and other important information, read the current plan
at specs/001-author-category/plan.md
<!-- SPECKIT END -->
