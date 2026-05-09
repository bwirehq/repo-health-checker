# Architecture

Repo Health Checker is organized around a reusable scanner engine. The CLI is intentionally thin.

## Package Responsibilities

- `cmd/repo-health`: process entrypoint.
- `internal/cli`: command definitions, flags, input validation, exit codes.
- `internal/github`: GitHub API access and repository data collection.
- `internal/local`: offline local repository data collection from files and git metadata.
- `internal/scanner`: orchestration layer that fetches repository data and runs checks.
- `internal/checks`: independent health checks.
- `internal/scoring`: score aggregation and grade calculation.
- `internal/report`: terminal and JSON rendering.
- `internal/config`: default thresholds and runtime settings.
- `internal/model`: shared domain types.

## Data Flow

1. The CLI parses input into a `RepoRef`.
2. The GitHub or local source fetches repository metadata, files, releases or tags, and recent commits.
3. The scanner passes the collected `RepositoryData` to each check.
4. Checks return typed `CheckResult` values.
5. The scoring package aggregates the results into a `ScanResult`.
6. The report package renders either text or JSON.

## Adding A Check

1. Add a type in `internal/checks`.
2. Implement `Run(context.Context, model.RepositoryData, config.Config) model.CheckResult`.
3. Add the check to `DefaultSuite`.
4. Add unit tests covering pass, warn, and fail behavior.
5. Document the scoring rule in `docs/scoring.md`.

Checks should not fetch data, print output, read environment variables, or know about CLI flags.

## Testing Strategy

Normal tests should be deterministic and avoid live network calls. GitHub API behavior should be tested through fixtures or small interfaces. Optional integration tests can be added behind `RUN_INTEGRATION=1`.
