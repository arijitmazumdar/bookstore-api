# Agent Customizations

## code-reviewer

- File: `.github/agents/code-reviewer.agent.md`
- Prompt: `.github/prompts/code-reviewer.prompt.md`

Use this agent to review the Bookstore API Go implementation and REST API design.
It is intended to help AI assistants deliver consistent, reusable review guidance across the repository.

### What it covers

- Go idioms and code quality
- REST API handler structure and routing semantics
- SQLite schema and SQL usage
- JSON serialization and HTTP error handling
- Test coverage and CI workflow validation
- Static analysis with `golangci-lint`
