# Data Model: Author Category

## Author

**Purpose**: Represents a writer in the bookstore catalog.

**Fields**:
- `id`: Unique identifier.
- `name`: Author display name.
- `category`: Derived classification returned by the API. Allowed values: `regular`, `premium`.
- `sold_copies`: Derived aggregate count of purchase records across all books for the author. This supports category calculation and test assertions, even if it is not exposed directly in every response.

**Relationships**:
- One author can have many books.
- One author can have many purchase records indirectly through books.

**Validation rules**:
- `name` is required for create and update operations.
- `category` is read-only from the client's perspective.
- `category = premium` only when `sold_copies > 500`; otherwise `category = regular`.

## Book

**Purpose**: Represents a catalog item written by one author.

**Fields used by this feature**:
- `id`
- `author_id`
- `name`

**Relationships**:
- Each book belongs to exactly one author.
- Each book can have many purchases.

## Purchase

**Purpose**: Represents one sold copy of a book.

**Fields used by this feature**:
- `id`
- `book_id`
- `customer_id`
- `purchase_date`
- `purchase_price`

**Relationships**:
- Each purchase belongs to exactly one book.
- Purchases contribute one sold copy each toward the linked author's `sold_copies` total.

## Derived State Rules

- `sold_copies` is the count of purchase rows joined to books owned by the author.
- Authors with no books have `sold_copies = 0`.
- Authors with books but no purchases also have `sold_copies = 0`.
- Category transitions are automatic:
  - `regular` -> `premium` when `sold_copies` becomes 501 or greater
  - `premium` -> `regular` when `sold_copies` falls back to 500 or lower
