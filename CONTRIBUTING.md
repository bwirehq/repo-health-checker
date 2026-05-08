# Contributing

Thanks for improving Repo Health Checker.

## Local Setup

```sh
go mod download
go test ./...
go build -o bin/repo-health ./cmd/repo-health
```

## Quality Checks

```sh
make fmt
make test
make vet
make lint
make check
```

`make lint` requires `golangci-lint`.

## Pull Request Expectations

- Keep package boundaries intact.
- Add tests for scoring or report changes.
- Avoid live GitHub calls in normal unit tests.
- Update docs when changing scoring behavior.
- Keep user-facing output stable unless the change is intentional.

## Design Principles

- Prefer clear types over maps.
- Prefer small interfaces over global state.
- Keep checks deterministic and independently testable.
- Make scores explainable.
