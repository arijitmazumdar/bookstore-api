# Feature Specification: Author Category

**Feature Branch**: `001-author-category`

**Created**: 2026-06-21

**Status**: Draft

**Input**: User description: "I want add author category as premium and regular. any books that are sold more than 500 copies will make the author as premium otherwise normal"

## Clarifications

### Session 2026-06-21

- Q: Which author API responses should include the derived `category` field? → A: `GET /authors`, `GET /authors/{id}`, and successful `POST`/`PUT` author responses all include `category`.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View author category in author responses (Priority: P1)

As a bookstore operator, I want each author record to show whether the author is `regular` or `premium` so that I can quickly identify high-performing authors.

**Why this priority**: The feature has no user value unless author category is visible in the existing author API.

**Independent Test**: Create an author with no purchases, request the author by ID and from the author list, and confirm the category is returned as `regular`.

**Acceptance Scenarios**:

1. **Given** an author has no recorded sales, **When** a client requests that author, **Then** the response includes category `regular`.
2. **Given** multiple authors exist with different sales totals, **When** a client requests the author list, **Then** every author entry includes a category derived from that author's current sales total.

---

### User Story 2 - Automatically upgrade an author to premium after enough sales (Priority: P2)

As a bookstore operator, I want an author's category to change to `premium` automatically once their books have sold more than 500 copies so that the system reflects current performance without manual updates.

**Why this priority**: The main business rule is automatic classification based on sales volume.

**Independent Test**: Create an author and a book, record 501 purchases for that author's books, then request the author and confirm the category is `premium`.

**Acceptance Scenarios**:

1. **Given** an author's books have a combined sales count of 501, **When** the author record is retrieved, **Then** the response includes category `premium`.
2. **Given** an author's books have a combined sales count of 500, **When** the author record is retrieved, **Then** the response still includes category `regular`.

---

### User Story 3 - Recalculate category when sales data changes (Priority: P3)

As a bookstore operator, I want author category to always reflect the latest purchase data so that reports do not become stale after purchases or catalog changes.

**Why this priority**: Derived business rules are only trustworthy if the category remains in sync with the underlying purchase data.

**Independent Test**: Change the set of purchases associated with an author's books and confirm later author reads reflect the updated category.

**Acceptance Scenarios**:

1. **Given** an author is `premium`, **When** related sales data changes and the total drops to 500 or below, **Then** the next author read shows category `regular`.

### Edge Cases

- Authors with no books must still be returned with category `regular`.
- Authors with multiple books must be categorized using the combined sales total across all of their books.
- A sales total of exactly 500 must remain `regular`; only totals greater than 500 qualify as `premium`.
- Invalid or failed purchase creation must not change any author category.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST return an author category for every successful author response exposed by the author API, including author list, author detail, author creation, and author update responses.
- **FR-002**: System MUST classify an author as `premium` when the combined sold-copy count across all books by that author is greater than 500.
- **FR-003**: System MUST classify an author as `regular` when the combined sold-copy count across all books by that author is 500 or fewer.
- **FR-004**: System MUST calculate sold-copy count from recorded purchase entries associated with books written by the author.
- **FR-005**: System MUST apply the same category rule consistently to author list responses and single-author responses.
- **FR-006**: System MUST treat author category as system-derived data rather than client-supplied input.
- **FR-007**: System MUST keep existing author creation and update workflows usable without requiring clients to send a category field.

### Key Entities *(include if feature involves data)*

- **Author**: A writer record identified by name and an automatically derived category based on current sales volume.
- **Book**: A catalog item linked to one author and contributing to that author's sales total through purchases.
- **Purchase**: A recorded sale of one book to one customer; each purchase contributes one sold copy toward the linked author's total.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of successful author list, author detail, author creation, and author update responses include a category value of either `regular` or `premium`.
- **SC-002**: An author with 501 recorded purchases across their books is returned as `premium` in validation tests.
- **SC-003**: An author with 500 or fewer recorded purchases across their books is returned as `regular` in validation tests.
- **SC-004**: Existing author creation and update API requests continue to succeed without clients providing a category field.

## Assumptions

- The user's word `normal` is interpreted as the `regular` category named in the request.
- One purchase record represents one sold copy for the associated book.
- Sales total is derived from current purchase records rather than entered manually by staff.
- This feature only changes author-facing API output and related validation coverage; separate reporting endpoints are out of scope.
