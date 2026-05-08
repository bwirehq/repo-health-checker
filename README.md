# Repo Health Checker

Repo Health Checker is a Go CLI that scans a public GitHub repository and reports a transparent health score. It is designed as a reusable scanner engine first, with the CLI as the first interface.

## Install

```sh
go install github.com/bwirehq/repo-health-checker/cmd/repo-health@latest
```

For local development:

```sh
go run . scan
go build -o bin/repo-health .
```

## Quick Start

```sh
go run . scan
repo-health scan
repo-health scan openai/codex
repo-health scan https://github.com/openai/codex --verbose
repo-health scan openai/codex --json
```

If you run `repo-health scan` without an argument, the CLI prompts for a repository:

```txt
GitHub repository (owner/repo or URL): https://github.com/openai/codex
```

Example output:

```txt
Repo Health: 78/100 (C)

PASS CI configured
  Continuous integration configuration was detected.
PASS Commit activity
  14 commits were found in the last 90 days.
WARN Issue health
  43 open issues; 19 appear stale.
FAIL Releases
  No releases or tags were found.

Top fixes:
1. Create a release - Tag a version and publish release notes when the project reaches a usable milestone.
```

## Scoring Rubric

The score is a transparent 100-point rubric:

| Category | Points |
| --- | ---: |
| README quality | 15 |
| Commit activity | 15 |
| Issue health | 15 |
| Pull request health | 10 |
| CI presence | 10 |
| License | 10 |
| Releases/tags | 10 |
| Test coverage hints | 10 |
| Dependency hygiene hints | 5 |

See [docs/scoring.md](docs/scoring.md) for detailed thresholds.

## GitHub Token

Public repositories work without authentication, but GitHub rate limits unauthenticated requests. Set `GITHUB_TOKEN` to raise the limit:

```sh
GITHUB_TOKEN=ghp_example repo-health scan owner/repo
```

The token is read from the environment only. It is never printed or written to disk.

## Flags

```txt
--json            write machine-readable JSON
--verbose         include score details in text output
--no-color        disable ANSI color output
--fail-under N    exit with code 2 when score is below N
--timeout 15s     GitHub API timeout
```

## Limitations

- V1 supports public GitHub repositories only.
- Dependency freshness is a lightweight hygiene signal, not a vulnerability audit.
- Test coverage is inferred from repository structure, not coverage reports.
- Normal tests use mocks and fixtures; live GitHub integration tests should be explicitly gated.

## Roadmap

- GitHub Action mode.
- Web dashboard using the same scanner engine.
- Organization-level reporting.
- Deeper dependency freshness checks.
- Historical score tracking.
