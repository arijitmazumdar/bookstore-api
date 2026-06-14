---
name: code-reviewer
description: |
  Review the Go bookstore API code for best practices, REST API design, and maintainability.
  Perform static code analysis with golangci-lint and report issues with actionable fixes.
applyTo:
  - "**/*.go"
  - ".github/workflows/*.yml"
---

# Code Review Prompt

Use this prompt to review the bookstore API source.

Tasks:
- Analyze Go code quality and idiomatic usage.
- Review REST API route design and JSON handling.
- Check database schema and SQL usage for correctness.
- Look for potential bugs, security issues, and maintainability problems.
- Run `golangci-lint run ./...` if available and summarize the findings.
- Provide specific, actionable recommendations.

If `golangci-lint` is not installed locally, mention that the review could be improved by running it in CI or a dev container.
