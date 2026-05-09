# Scoring

Repo Health Checker uses a transparent 100-point weighted rubric. The goal is to identify maintenance signals that developers can understand and improve.

| Category | Points |
| --- | ---: |
| README quality | 10 |
| CI presence | 15 |
| License | 10 |
| Commit activity | 20 |
| Issue health | 10 |
| Pull request health | 10 |
| Releases/tags | 10 |
| Test coverage hints | 10 |
| Dependency hygiene hints | 5 |

Archived repositories receive a 10-point penalty after the weighted score is calculated. They are always reported as high risk because GitHub marks them read-only.

## README Quality: 10 points

- Pass: README is substantial and includes structure plus setup or usage signals.
- Warn: README exists but is thin or missing setup, usage, testing, or contribution details.
- Fail: no README detected.

## CI Presence: 15 points

- Pass: GitHub Actions workflows or common CI config files are detected.
- Fail: no common CI configuration is detected.

## License: 10 points

- Pass: GitHub license metadata or a common license file is detected.
- Fail: no license is detected.

## Commit Activity: 20 points

- Pass: at least 10 commits in the last 90 days.
- Warn: at least 1 commit in the last 90 days.
- Fail: no commits in the last 90 days.

## Issue Health: 10 points

Issues are considered stale when they have not been updated for 90 days.

- Pass: no open issues, or no more than 20 percent are stale.
- Warn: more than 20 percent and no more than 50 percent are stale.
- Fail: more than 50 percent are stale.
- Info: unavailable during offline local scans.

## Pull Request Health: 10 points

Pull requests are considered stale when they have not been updated for 30 days.

- Pass: no open pull requests, or none are stale.
- Warn: 1 to 3 stale pull requests.
- Fail: more than 3 stale pull requests.
- Info: unavailable during offline local scans.

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

## Risk Levels

- Low: score is at least 80 and no check is warning or failing.
- Medium: score is at least 60, or at least one non-critical check is warning or failing.
- High: repository is archived, score is below 60, CI is missing, or commit activity fails.
