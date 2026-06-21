# Research: Author Category

## Decision 1: Derive author category from purchase data instead of storing it manually

**Decision**: Compute author category from current purchase records linked through the author's books rather than persisting a manually edited category value.

**Rationale**: The requested rule is based entirely on sales totals. Deriving category from existing purchase data avoids duplication, prevents drift between stored status and actual sales, and keeps the source of truth in one place.

**Alternatives considered**:
- Store `category` directly on the `authors` table and update it during purchase creation.
- Store cumulative sold-copy counts on authors and recompute category from that value.

## Decision 2: Count sold copies as the number of purchase rows for an author's books

**Decision**: Treat each row in `customer_book_purchase` as one sold copy and aggregate counts across every book associated with the author.

**Rationale**: The current data model records one purchase per book per row and does not include a quantity field. This gives a clear, testable rule without changing the purchase schema.

**Alternatives considered**:
- Interpret `purchase_price` as a quantity signal.
- Add a new quantity field to purchases before implementing author category.

## Decision 3: Expose category on read responses and keep write payloads unchanged

**Decision**: Return `category` in author responses while continuing to accept existing author create and update payloads without requiring category input.

**Rationale**: The category is system-derived, so clients should not provide or control it. Keeping write payloads stable minimizes client impact and matches the requirement that category is computed from sales.

**Alternatives considered**:
- Require clients to submit a category when creating or updating an author.
- Create a dedicated endpoint for author category rather than extending author responses.

## Decision 4: Use the threshold rule literally: premium only when sales are greater than 500

**Decision**: Authors become `premium` only when total sold copies are strictly greater than 500. At exactly 500 they remain `regular`.

**Rationale**: The request explicitly says "more than 500 copies." Preserving that exact rule removes ambiguity and gives a crisp test boundary.

**Alternatives considered**:
- Promote to `premium` at 500 or more.
- Introduce additional tiers beyond `regular` and `premium`.
