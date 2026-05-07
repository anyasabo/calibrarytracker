# Agent Notes

Guidance for coding agents working in this repository.

## Project Shape

- Frontend: SvelteKit static app in `src/`
- Scraper: Go CLI in `cmd/scraper/` with implementation in `internal/scraper/`
- Reference data: checked in under `data/` and consumed by frontend

## Runtime and Version Sources

Use the repo version files as source of truth:

- Node: `.node-version`
- Go: `.go-version`

Keep these in sync with:

- `package.json` -> `engines.node`
- `go.mod` -> `go` directive

CI enforces these alignments.

## Local Workflow

Preferred setup:

```bash
mise install
npm ci --no-audit --no-fund
```

If `mise` fails to install Node because of a literal `24.x` entry, use mise-compatible syntax in `.tool-versions` (`node prefix:24`).

### Command Execution Policy (Agents)

Use `mise exec --` for toolchain-dependent commands so shells do not need `mise activate` preloaded.

Examples:

```bash
mise exec -- npm ci --no-audit --no-fund
mise exec -- npm run lint
mise exec -- npm run check
mise exec -- npm run build
mise exec -- npm run test:smoke
mise exec -- go vet ./...
mise exec -- go test ./...
```

## Required Checks Before Pushing

Run the relevant checks for changed areas:

```bash
npm run lint
npm run check
npm run build
go vet ./...
go test ./...
```

## CI/Dependency Guardrails

- Keep `package-lock.json` registry URLs on `registry.npmjs.org` (do not reintroduce private Artifactory URLs), otherwise GitHub-hosted CI/Dependabot can time out.
- `ci.yml` uses path-based gating; if you add new source/config paths that affect frontend or scraper, update the filters.
- `deploy.yml` deploys from the CI artifact via `workflow_run` and supports manual `workflow_dispatch`.

## Data Rules

- Do not hand-edit `data/systems.json` or `data/branches.json`; regenerate via scraper.
- Hand-edited reference files (`data/cooperatives.json`, `data/partnerships.json`, and `data/digital-access.json`) may be updated directly.
