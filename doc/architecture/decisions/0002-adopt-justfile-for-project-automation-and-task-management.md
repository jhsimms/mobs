# 2. Adopt justfile for project automation and task management

Date: 2025-05-08

## Status

Accepted

## Context

I needed a simple, consistent way to automate tasks like build, test, lint, and deploy. Bash scripts were scattered and inconsistent. I wanted something cross-platform, easy to read, and lightweight.

## Decision

I adopted just as the task runner. The justfile now defines all common workflows. It’s the single source of truth for project automation.

## Consequences
- Tasks are discoverable via just --list
- Logic is centralized and readable
- Cross-platform and easy to extend
- Requires installing just
- Slight learning curve for new contributors

## Alternatives
- Makefile: Complex, not great on Windows
- Bash: Flexible but messy
- Taskfile: Feature-rich but heavier

just hit the right balance of simplicity and power.