# Quickstart: Author Category

## Purpose

Validate that author responses expose a derived `category` and that the `premium` threshold is applied correctly from purchase data.

## Prerequisites

- Go 1.22 installed
- Writable temporary directory for SQLite test databases
- Repository dependencies installed through normal Go module resolution

## Validation Scenarios

### Scenario 1: New author defaults to regular

1. Run the targeted handler or integration test covering author creation and retrieval.
2. Create an author with no books or purchases.
3. Retrieve that author through the API or inspect the successful create/update response payload.

**Expected outcome**: Every successful author response (`POST`, `PUT`, `GET /authors`, and `GET /authors/{id}`) includes `category: "regular"` when the author has no qualifying sales.

### Scenario 2: Author becomes premium after more than 500 sales

1. Create one author and at least one book for that author.
2. Record 501 purchases for the author's books.
3. Retrieve the author through the API.

**Expected outcome**: The author response includes `category: "premium"`.

### Scenario 3: Boundary value of exactly 500 remains regular

1. Create one author and at least one book for that author.
2. Record exactly 500 purchases across that author's books.
3. Retrieve the author through the API.

**Expected outcome**: The author response includes `category: "regular"`.

### Scenario 4: Category recalculates from current stored purchase data

1. Create one author and at least one book for that author.
2. Record 501 purchases for the author's books and verify the author is `premium`.
3. Remove enough related purchase rows in a test fixture so the total drops to 500 or below.
4. Retrieve the author again through the API.

**Expected outcome**: The next author response includes `category: "regular"`.

## Suggested Commands

```bash
go test ./internal/handlers ./internal/app
go test ./... 
```

## Contract Reference

- `specs/001-author-category/contracts/authors.yaml`

## Data Model Reference

- `specs/001-author-category/data-model.md`
