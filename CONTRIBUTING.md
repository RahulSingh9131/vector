# Contributing to Vector

Thank you for your interest in contributing to Vector! This guide will help you get started.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Commit Convention](#commit-convention)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

Please be respectful and constructive in all interactions. We are committed to providing a welcoming and inclusive environment for everyone.

## Getting Started

1. **Fork** the repository on GitHub
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/<your-username>/vector.git
   cd vector
   ```
3. **Set up your environment** — follow the [Local Development Setup](README.md#-quick-start) in the README
4. **Create a branch** for your work:
   ```bash
   git checkout -b feat/your-feature-name
   ```

## Development Workflow

### Backend (Go)

```bash
cd backend

# Run the server
task run

# Run the linter
task lint

# Auto-fix lint issues
task lint:fix

# Format and tidy dependencies
task tidy
```

### Frontend Packages (TypeScript)

```bash
# Install dependencies
bun install

# Build all packages
bun run build

# Lint all packages
bun run lint

# Type-check all packages
bun run type-check
```

### Database Migrations

```bash
cd backend

# Create a new migration
task migrations:new name=your_migration_name

# Apply migrations
task migrations:up
```

## Commit Convention

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification. Each commit message should be structured as follows:

```
<type>(<scope>): <description>

[optional body]
```

### Types

| Type       | Description                                 |
|------------|---------------------------------------------|
| `feat`     | A new feature                               |
| `fix`      | A bug fix                                   |
| `docs`     | Documentation only changes                  |
| `style`    | Code style changes (formatting, semicolons) |
| `refactor` | Code change that neither fixes nor adds     |
| `perf`     | Performance improvement                     |
| `test`     | Adding or correcting tests                  |
| `chore`    | Maintenance tasks (deps, configs)           |

### Examples

```
feat(issues): add drag-and-drop reordering
fix(auth): handle expired session tokens
docs: update README with setup instructions
refactor(database): extract connection pooling logic
```

## Pull Request Process

1. **Ensure your branch is up to date** with `main`:
   ```bash
   git fetch origin
   git rebase origin/main
   ```
2. **Run all checks locally** before opening a PR:
   ```bash
   # Backend
   cd backend && task lint && task tidy

   # Frontend
   bun run lint && bun run type-check && bun run build
   ```
3. **Open a Pull Request** with:
   - A clear, descriptive title following commit conventions
   - A summary of what changed and why
   - Screenshots for any UI changes
   - Links to related issues (e.g., `Closes #123`)
4. **Address review feedback** promptly and push updates to the same branch

## Coding Standards

### Go (Backend)

- Follow the [Effective Go](https://go.dev/doc/effective-go) guidelines
- Every package must have a package-level doc comment (`// Package xyz ...`)
- Use structured logging via `zerolog` — avoid `fmt.Println` or `log`
- All errors must be wrapped with context: `fmt.Errorf("doing X: %w", err)`
- Keep handlers thin — business logic belongs in the **service layer**
- Database access goes through the **repository layer**
- Use the `DBTX` interface for transactional operations
- Run `task lint` before committing

### TypeScript (Frontend Packages)

- Use strict TypeScript — no `any` unless absolutely necessary
- Keep shared contracts in the `@vector/openapi` package
- Validation schemas belong in the `@vector/zod` package

### Project Structure

- **Handlers** → Parse requests, call services, return responses
- **Services** → Business logic, authorization, orchestration
- **Repositories** → Database queries (single responsibility)
- **Events** → Pub/sub for activity tracking and side effects
- **Models** → Domain types shared across layers

## Reporting Issues

When opening an issue, please include:

1. **Description** — what happened vs. what you expected
2. **Steps to reproduce** — minimal steps to trigger the issue
3. **Environment** — OS, Go version, Node version, Docker version
4. **Logs / Screenshots** — any relevant output or visuals

---

Thank you for contributing to Vector! 🚀
