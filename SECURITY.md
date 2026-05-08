# Security Policy

## Supported Versions

Security fixes are applied to the latest released version.

## Token Handling

Repo Health Checker reads `GITHUB_TOKEN` from the environment when present. Tokens are not logged, printed, persisted, or included in JSON output.

## Reporting A Vulnerability

Please report security issues privately to the maintainers. Include:

- affected version or commit
- impact
- reproduction steps
- any relevant logs with secrets removed

## Scope

V1 scans public GitHub repositories through GitHub APIs. It does not execute repository code, clone repositories, or run package manager commands.
