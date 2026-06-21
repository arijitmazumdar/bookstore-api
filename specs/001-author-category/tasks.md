# Tasks: Author Category

**Input**: Design documents from `/specs/001-author-category/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Include targeted handler and integration tests because the repository guidance requires updating tests when API responses and SQL behavior change.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Go API source lives under `internal/`
- Integration coverage lives in `internal/app/`
- Handler-level coverage lives in `internal/handlers/`
- Feature design artifacts live in `specs/001-author-category/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm the feature artifacts and validation targets are aligned before implementation

- [X] T001 Review author-category contract and validation scenarios in `specs/001-author-category/contracts/authors.yaml` and `specs/001-author-category/quickstart.md`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared author response structures and sales aggregation support that every story depends on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 Add derived author category fields and response shape updates in `internal/models/models.go`
- [X] T003 Implement shared author sales aggregation query helpers in `internal/handlers/author.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in priority order

---

## Phase 3: User Story 1 - View author category in author responses (Priority: P1) 🎯 MVP

**Goal**: Return `regular` or `premium` in author list and author detail responses without changing author write payload requirements

**Independent Test**: Create an author with no purchases, then verify `/authors` and `/authors/{id}` both return `category: "regular"`.

### Tests for User Story 1

- [X] T004 [P] [US1] Add handler coverage for regular-category author list/detail responses in `internal/handlers/handlers_test.go`
- [X] T005 [P] [US1] Extend end-to-end author response assertions for derived category in `internal/app/integration_test.go`

### Implementation for User Story 1

- [X] T006 [US1] Return derived category from author list and author detail handlers in `internal/handlers/author.go`
- [X] T007 [US1] Preserve author create and update request compatibility while returning derived category in `internal/handlers/author.go`

**Checkpoint**: User Story 1 should now be independently functional and testable as the MVP

---

## Phase 4: User Story 2 - Automatically upgrade an author to premium after enough sales (Priority: P2)

**Goal**: Mark authors as `premium` once aggregate purchase count across their books exceeds 500, while keeping exactly 500 as `regular`

**Independent Test**: Create one author and one book, record 500 purchases and verify `regular`, then add the 501st purchase and verify `premium`.

### Tests for User Story 2

- [X] T008 [P] [US2] Add handler coverage for the 500-versus-501 sales threshold in `internal/handlers/handlers_test.go`
- [X] T009 [P] [US2] Extend integration coverage for premium-category threshold behavior in `internal/app/integration_test.go`

### Implementation for User Story 2

- [X] T010 [US2] Apply the premium threshold rule to aggregated sold-copy counts in `internal/handlers/author.go`
- [X] T011 [US2] Ensure author reads aggregate purchases across all books by the same author in `internal/handlers/author.go`

**Checkpoint**: User Stories 1 and 2 should both work independently, including the exact threshold boundary

---

## Phase 5: User Story 3 - Recalculate category when sales data changes (Priority: P3)

**Goal**: Keep category derived from current purchase data so later reads reflect any increase or decrease in qualifying sales volume

**Independent Test**: Seed an author above the premium threshold, remove enough related purchases in a test fixture, and verify the next author read returns `regular`.

### Tests for User Story 3

- [X] T012 [P] [US3] Add handler coverage for category recalculation after purchase removal or data reset in `internal/handlers/handlers_test.go`

### Implementation for User Story 3

- [X] T013 [US3] Rework author category reads to derive status from current joined purchase data on every request in `internal/handlers/author.go`

**Checkpoint**: All user stories should now be independently functional with category recalculated from live purchase data

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and cleanup across all stories

- [X] T014 [P] Update author API examples and feature notes in `specs/001-author-category/quickstart.md` if implementation details require clarification
- [X] T015 Run full regression validation for the feature with `go test ./...`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational completion
- **User Story 2 (Phase 4)**: Depends on Foundational completion and builds naturally on User Story 1 response work
- **User Story 3 (Phase 5)**: Depends on Foundational completion and on the threshold logic from User Story 2
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational - no dependency on other stories
- **User Story 2 (P2)**: Depends on User Story 1 author response coverage so the threshold behavior has a visible output surface
- **User Story 3 (P3)**: Depends on User Story 2 because recalculation uses the same derived threshold rule

### Within Each User Story

- Tests should be added before or alongside implementation and must fail until the story behavior is implemented
- Shared model and aggregation changes come before handler response updates
- Handler behavior must be in place before full integration assertions are considered complete

### Parallel Opportunities

- `T004` and `T005` can run in parallel after foundational work
- `T008` and `T009` can run in parallel after User Story 1 is complete
- `T012` can be prepared in parallel with implementation planning for `T013`
- `T014` and `T015` can run in parallel once implementation is complete, assuming documentation updates do not affect code

---

## Parallel Example: User Story 1

```bash
# Launch the story-level validation tasks together after foundational work:
Task: "Add handler coverage for regular-category author list/detail responses in internal/handlers/handlers_test.go"
Task: "Extend end-to-end author response assertions for derived category in internal/app/integration_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. Stop and validate `/authors` and `/authors/{id}` return derived categories for low-sales authors

### Incremental Delivery

1. Finish the shared model and aggregation support
2. Deliver User Story 1 for visible category output
3. Add User Story 2 to enforce the premium threshold
4. Add User Story 3 to prove recalculation stays live with purchase data changes
5. Run the full regression suite

### Parallel Team Strategy

With multiple developers:

1. One developer handles foundational aggregation support
2. A second developer prepares handler and integration tests for the next ready story
3. Once User Story 1 lands, threshold and recalculation test work can be split across stories without touching unrelated files

---

## Notes

- [P] tasks target different files or documentation and can be done without blocking each other once dependencies are met
- Each user story phase is designed to be independently demonstrable
- The suggested MVP scope is User Story 1 only
- All tasks use the required checklist format with task ID, optional parallel marker, user story label where required, and explicit file path
