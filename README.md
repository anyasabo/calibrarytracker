# CA Library Card Tracker

A static web app for tracking California public library cards, branch visits, and reciprocal borrowing access.

## What This Repo Contains

- `src/`: SvelteKit frontend (static site)
- `cmd/scraper/` + `internal/scraper/`: Go scraper that refreshes reference data
- `data/`: checked-in library data files used by the app

## Quick Start

```bash
# install runtimes from repo version files
mise install

# install frontend dependencies
npm ci --no-audit --no-fund

# run the app locally
npm run dev
```

## Common Commands

```bash
# refresh library reference data
go run ./cmd/scraper

# frontend checks/build
npm run lint
npm run check
npm run build

# scraper checks
go vet ./...
go test ./...
```

## Deployment

- GitHub Actions runs CI checks on pushes/PRs.
- GitHub Pages deploys from successful `main` CI output.

## Contributing

See `CONTRIBUTING.md` for full workflow details and project conventions.
