# 2. Adopt justfile for project automation and task management

Date: 2025-05-08

## Status

Accepted

## Context

As the project grew, we needed a consistent, maintainable, and developer-friendly way to automate common tasks such as building, testing, linting, formatting, and managing infrastructure. Previously, these tasks were handled via ad-hoc Bash scripts or manual commands, which led to inconsistencies, onboarding friction, and duplicated logic. We wanted a solution that would:

- Be easy to read and modify
- Work cross-platform (Linux, macOS, Windows via WSL)
- Support parameterized tasks and dependencies
- Be lightweight and require minimal setup
- Encourage documentation and discoverability of available tasks

## Decision

We adopted a justfile (using the `just` task runner) as the canonical way to define and run project automation tasks. The justfile now serves as the single source of truth for build, test, lint, format, infrastructure, and other developer workflows. All contributors are expected to use and update the justfile for automation needs.

## Consequences

- **Easier onboarding:** New contributors can discover and run all project tasks via `just --list`.
- **Consistency:** All automation logic is centralized, reducing duplication and drift.
- **Readability:** The justfile syntax is simple and self-documenting, making it easy to understand and modify tasks.
- **Cross-platform support:** just works on all major platforms, reducing environment-specific issues.
- **Extensibility:** Adding new tasks or parameters is straightforward.
- **Dependency:** Contributors must install `just` (a small Rust binary), which is an extra step compared to using only Bash or Make.
- **Learning curve:** Some team members may be unfamiliar with justfile syntax, though it is minimal.

## Alternatives

- **Makefile:** Widely used, but Make has a more complex syntax, poor cross-platform support (especially on Windows), and is less user-friendly for non-C/C++ projects.
- **Bash scripts:** Flexible and ubiquitous, but can become unwieldy, hard to maintain, and lack built-in task discovery/documentation.
- **Taskfile (go-task):** More features, but heavier and requires Go; just is lighter and simpler for our needs.
- **npm scripts:** Good for JS projects, but not suitable for a polyglot stack or non-JS contributors.

We chose justfile because it best balanced simplicity, readability, cross-platform support, and ease of use for our team and project needs.
