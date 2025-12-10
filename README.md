# Vector

A modern monorepo powered by [Bun](https://bun.sh) and [Turbo](https://turbo.build).

## Prerequisites

- [Bun](https://bun.sh) v1.3.4 or higher
- Node.js >= 18.0.0

## Getting Started

### Installation

Install dependencies using Bun:

```bash
bun install
```

### Development

Run all packages in development mode:

```bash
bun run dev
```

### Building

Build all packages:

```bash
bun run build
```

## Available Scripts

| Script | Description |
|--------|-------------|
| `bun run dev` | Start development mode for all packages |
| `bun run build` | Build all packages |
| `bun run lint` | Lint all packages |
| `bun run test` | Run tests across all packages |
| `bun run type-check` | Run TypeScript type checking |
| `bun run format` | Format code across all packages |
| `bun run clean` | Clean build artifacts |
| `bun run turbo:clean` | Deep clean (removes .turbo cache and node_modules) |

## Project Structure

```
vector/
├── packages/          # Workspace packages
├── turbo.json        # Turbo configuration
├── package.json      # Root package configuration
└── README.md         # This file
```

## Workspaces

This project uses Bun workspaces. Add new packages in the `packages/` directory.

## Tech Stack

- **Package Manager**: Bun
- **Build System**: Turbo
- **Language**: TypeScript

## License

ISC
