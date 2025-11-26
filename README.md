# Learning Cards

A small Go-based backend for a spaced-repetition "learning cards" service. The project provides:
- A REST API to fetch daily/review words and words by category.
- A lightweight spaced-repetition scheduling system for user word reviews.
- Automatic seeding of words from CSV files in the `data/` directory.
- Postgres-backed persistence via Gorm.
- Cron-based background sync and seeding (different schedules for local vs production).

## Table of Contents

- [Features](#features)
- [Tech stack](#tech-stack)
- [Repository layout](#repository-layout)
- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Run locally](#run-locally)
- [Run with Docker / docker-compose](#run-with-docker--docker-compose)
- [API Reference](#api-reference)
- [Data / CSVs](#data--csvs)
- [Cron behavior](#cron-behavior)
- [License](#license)

## Features

- Retrieve user words due for review today.
- Retrieve user words by category (only those due for review).
- Update a word's learning status (learned / failed) and update scheduling.
- Seed words from CSV files in `data/`.
- Automatic DB migrations on startup.

## Tech stack

- Go 1.25
- Gin (HTTP framework)
- Gorm (ORM) with the `postgres` driver
- Postgres database
- Robfig/cron for scheduled jobs

## Repository layout (key files)

- `internal/cmd/main.go` — application entrypoint
- `internal/startup/startup.go` — app wiring, router & cron setup
- `api/v1/routes.go` — route registration
- `internal/handlers` — HTTP handlers
- `internal/services` — business logic
- `internal/repository` — DB access (Gorm)
- `internal/models` — Gorm models
- `internal/database/db.go` — DB connection
- `internal/database/migrations.go` — auto-migration logic
- `internal/utils` — CSV loader
- `data/` — CSV files used for seeding
- `config/config.go` — environment-based configuration (includes `DefaultDBConfig`)

## Prerequisites

- Go 1.25 or later installed
- Postgres (local or remote) for persistence
- Docker & docker-compose (optional, for running via containers)

## Configuration

The application reads configuration from environment variables. Important variables:

- `DB_HOST` — Postgres host (default: `localhost`)
- `DB_USER` — Postgres user (default: `defaultuser`)
- `DB_PASSWORD` — Postgres password (default: `defaultpassword`)
- `DB_PORT` — Postgres port (default: `5432`)
- `DB_SSLMODE` — Postgres sslmode (default: `disable`)
- `DB_NAME` — Postgres database name (used by `docker-compose`, note: the DB name is optional in the app DSN depending on the environment)

You can create a `.env` file (not committed) and export these variables, or set them in your shell.

Example environment (fill with your values):
  - `DB_HOST=localhost`
  - `DB_USER=learning_user`
  - `DB_PASSWORD=supersecret`
  - `DB_PORT=5432`
  - `DB_SSLMODE=disable`
  - `DB_NAME=learning_cards`

The app provides `DefaultDBConfig` in `config/config.go` and uses `getEnv(key, default)` to fall back to defaults if env vars are missing.

## Run locally

1. Ensure Postgres is running and reachable with the env vars above.
2. From the repository root, run:

- Quick run (uses `go run`):
  - `go run ./internal/cmd`

- Or build & run:
  - `go build -o learning-cards .`
  - `./learning-cards`

On startup the app will:
- Open the database connection (see `internal/database/db.go`).
- Auto-migrate models (`internal/database/migrations.go`).
- Load seed data from the `data/` CSVs and attempt to insert missing words.
- Start an HTTP server (default port: `:8080`) and cron jobs.

If you prefer, set the environment and run via your IDE / debugger.

## Linting

This project uses `golangci-lint` for static analysis. The repository includes a pre-commit hook configuration (`.pre-commit-config.yaml`) to run the linter locally on commit. Below are recommended instructions for running lint checks locally, via pre-commit, and in CI.

### Local

1. Install Go (match the project's Go version — currently Go 1.25).
2. Install `golangci-lint` (build with your Go toolchain so the binary is compatible with your target Go version):
   - `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8`
   - Make sure `$(go env GOPATH)/bin` or `$HOME/go/bin` is on your `PATH`.
3. Run the linter:
   - `golangci-lint run --timeout=5m --out-format=colored-line-number ./...`
   - For quieter/targeted runs you can scope to packages, e.g. `./internal/...`.

Notes:
- The `github-actions` output format is deprecated; prefer `colored-line-number` or the default formats.
- golangci-lint can refuse to run if the linter binary was built with an older Go version than your project's `go` version — build/install it with the same Go toolchain used for development or in CI.

### Pre-commit

The repository includes a `.pre-commit-config.yaml` that defines a `golangci-lint` hook using `language: system`. To enable:

1. Install `pre-commit` (e.g. `pip install pre-commit`).
2. Ensure `golangci-lint` is installed and on your `PATH` (see Local steps).
3. Install the git hook in your clone:
   - `pre-commit install`
4. Optionally run lint on all files:
   - `pre-commit run --all-files`

The pre-commit hook runs at commit time and by default runs the linter across the module (not only staged files) to allow whole-package checks. If you prefer faster staged-file checks, consider configuring a lightweight set of checks for pre-commit and leave the full suite to CI.

### CI (GitHub Actions)

To ensure consistent linter builds in CI, build `golangci-lint` inside the job (so the binary is built with the job's Go version) or use the official action with an install mode that builds from source. Example steps that build and run the linter:

- Install with the job's Go toolchain:
```yaml
- name: Set up Go
  uses: actions/setup-go@v4
  with:
    go-version: '1.25'

- name: Install golangci-lint
  run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

- name: Run golangci-lint
  run: |
    golangci-lint --version
    golangci-lint run --timeout=5m --verbose --out-format=colored-line-number ./...
```

Alternatively, if using the `golangci/golangci-lint-action@v4`, prefer an install mode that builds with the job Go (for example `install-mode: goinstall`) or explicitly install the linter as above before running it — this prevents the "binary built with older Go than targeted" error.

CI should fail when linter errors are found (so code must be fixed before merging), but you can temporarily allow lint failures with `continue-on-error: true` on the step while addressing issues.

## Run with Docker / docker-compose

There is a `Dockerfile` and `docker-compose.yml` for development convenience.

- Build image locally:
  - `docker build -t learning-cards .`

- Run with docker-compose:
  - `docker compose up --build`

The provided `docker-compose.yml` starts two services:
- `db` — Postgres (exposes `5432`)
- `app` — your Go application (exposes `8080`)

Make sure your `.env` values are set or export the required environment variables before running `docker-compose`.

## Database / Migrations

This project uses Gorm's `AutoMigrate` to create/update tables for:
- `Word`
- `UserWord`

Migration is triggered on application startup by `database.Migrate(db)` in `internal/startup/startup.go`.

If you want more advanced migration support (versioned migrations), consider integrating a migration tool like `golang-migrate/migrate`.

## API Reference

All routes are registered under `/v1` (see `api/v1/routes.go`).

1. GET `/v1/words/daily`
   - Description: Returns the list of user words due today (shuffled).
   - Response: JSON array of `UserWord` objects (each preloads `Word`).
   - Example: `curl http://localhost:8080/v1/words/daily`

2. GET `/v1/words/category/:category`
   - Description: Returns user words due for review filtered by `category`.
   - Params:
     - `category` — category string defined in the CSV/words (e.g., `animals`, `food`)
   - Example: `curl http://localhost:8080/v1/words/category/animals`

3. PUT `/v1/words/update/:wordID`
   - Description: Update learning status for a word (mark as learned or not).
   - Params:
     - `wordID` — numeric ID of the word in `words` table or user words.
   - Body (JSON):
     - `{ "learned": true }` or `{ "learned": false }`
   - Example: `curl -X PUT -H "Content-Type: application/json" -d '{"learned":true}' http://localhost:8080/v1/words/update/123`

Responses use standard HTTP status codes. Errors generally return a JSON payload with an `error` field.

## Data / CSVs

Seed data is stored in `data/` as CSV files (example: `data/animals.csv`, `data/food.csv`, ...).

CSV format expectation (3 columns, header row supported):
- `word,translation,category`

The CSV loader:
- Reads all CSVs in the `data/` directory (`internal/utils/csvloader.go`).
- Converts records to `models.Word`.
- On cron or startup, the `insertData` job checks if a word exists and inserts it if missing.

If you add new CSVs, they will be processed on next run / cron seeding.

## Cron behavior

Cron jobs are configured in `internal/startup/cron.go`. Behavior:

- If running on a local machine (hostname equals `localhost` or matches the configured `HostnameIP`/`Hostname` in `config.LoadAppConfig()`), cron runs every minute for:
  - Syncing user words
  - Inserting CSV data (useful during development)

- In production mode (non-local hostnames), cron jobs run less frequently:
  - Sync user words: daily at 00:00
  - Insert CSV data: daily at 01:00

The cron jobs call:
- `handler.SyncUserWords()` — syncs any missing words between `words` and `user_words`.
- `insertData(db, words)` — attempts to insert predefined words from CSVs.

## Logging & errors

- The app logs to stdout using the standard library `log`.
- DB connection errors are returned from `database.Open()` and will cause the app not to start.
- Warnings are logged when CSV loading fails or DB ping fails.

## Next improvements / ideas

- Add versioned migrations (e.g., `migrate`).
- Add authentication & user management.
- Add comprehensive integration tests (HTTP + DB).
- Allow configuring cron schedules via environment or config file.
- Provide a simple frontend UI to review cards.

## License

This project is provided under the MIT License. See `LICENSE` for details (if you add one).

## Contact

For questions or help, open an issue or reach out via PR comments.
