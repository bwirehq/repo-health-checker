# Scoring

Repo Health Checker uses a transparent 100-point rubric. The goal is to identify maintenance signals that developers can understand and improve.

## README Quality: 15 points

- Pass: README is substantial and includes structure plus setup or usage signals.
- Warn: README exists but is thin or missing setup, usage, testing, or contribution details.
- Fail: no README detected.

## Commit Activity: 15 points

- Pass: at least 10 commits in the last 90 days.
- Warn: at least 1 commit in the last 90 days.
- Fail: no commits in the last 90 days.

## Issue Health: 15 points

Issues are considered stale when they have not been updated for 90 days.

- Pass: no open issues, or no more than 20 percent are stale.
- Warn: more than 20 percent and no more than 50 percent are stale.
- Fail: more than 50 percent are stale.

## Pull Request Health: 10 points

Pull requests are considered stale when they have not been updated for 30 days.

- Pass: no open pull requests, or none are stale.
- Warn: 1 to 3 stale pull requests.
- Fail: more than 3 stale pull requests.

## CI Presence: 10 points

- Pass: GitHub Actions workflows or common CI config files are detected.
- Fail: no common CI configuration is detected.

## License: 10 points

- Pass: GitHub license metadata or a common license file is detected.
- Fail: no license is detected.

## Releases And Tags: 10 points

- Pass: at least one GitHub release exists.
- Warn: tags exist but no GitHub releases exist.
- Fail: no releases or tags exist.

## Test Coverage Hints: 10 points

- Pass: common test files or test directories are detected.
- Warn: no common test files or test directories are detected.

This is a hint, not a coverage measurement.

## Dependency Hygiene: 5 points

- Pass: dependency manifests and common lockfiles are detected.
- Warn: dependency manifests exist but no common lockfile is detected.
- Info: no dependency manifests are detected.

This category intentionally avoids claiming vulnerability coverage in v1.
