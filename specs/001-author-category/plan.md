# Implementation Plan: Author Category

**Branch**: `001-author-category` | **Date**: 2026-06-21 | **Spec**: `specs/001-author-category/spec.md`

**Input**: Feature specification from `specs/001-author-category/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Add a derived `category` field to author API responses so authors become `premium` when the combined purchase count across all of their books exceeds 500 and remain `regular` otherwise. The implementation will compute category from existing `books` and `customer_book_purchase` data instead of introducing manual category management.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.22

**Primary Dependencies**: Go standard library `net/http`, `database/sql`, `encoding/json`, `github.com/mattn/go-sqlite3`

**Storage**: SQLite

**Testing**: Go `testing`, `httptest`, integration tests in `internal/app`, handler tests in `internal/handlers`

**Target Platform**: Linux-hosted REST API

**Project Type**: Web service

**Performance Goals**: Author list and detail reads continue to behave like current CRUD endpoints and remain suitable for local SQLite-backed test and development workloads

**Constraints**: Preserve existing author create/update request shape, keep migrations deterministic, maintain referential integrity, avoid manual category writes that can drift from sales data

**Scale/Scope**: Single-service change touching author responses, purchase aggregation logic, and automated tests for threshold behavior

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The current `.specify/memory/constitution.md` contains placeholder text only and does not define enforceable project-specific gates.

- Gate status before research: PASS
- Gate status after design: PASS
- Followed repo guidance from `AGENTS.md`: add focused test coverage, keep changes resource-oriented, preserve existing API behavior where possible

## Project Structure

### Documentation (this feature)

```text
specs/001-author-category/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/
│   └── authors.yaml     # Author API contract updates
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
cmd/server/main.go
internal/app/
├── routes.go
├── server.go
└── integration_test.go
internal/db/
├── db.go
└── migrations.go
internal/handlers/
├── author.go
├── book.go
├── purchase.go
└── handlers_test.go
internal/models/
└── models.go
scripts/
└── run-*.sh
```

**Structure Decision**: Keep the existing single-service Go API structure. The feature belongs in `internal/models`, `internal/handlers`, `internal/db`, and targeted tests under `internal/handlers` and `internal/app`.

## Complexity Tracking

No constitution violations or exceptional complexity justifications are required for this feature.
