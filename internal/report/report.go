package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/bwirehq/repo-health-checker/internal/model"
)

type Options struct {
	JSON     bool
	Compact  bool
	NoColor  bool
	Verbose  bool
	Duration time.Duration
}

func Write(w io.Writer, result model.ScanResult, opts Options) error {
	if opts.JSON {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}
	if opts.Compact {
		return writeCompact(w, result)
	}
	return writeText(w, result, opts)
}

func writeText(w io.Writer, result model.ScanResult, opts Options) error {
	_, err := fmt.Fprintf(w, "Repo Health: %d/%d (%s)\n", result.Score, result.MaxScore, result.Grade)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Repository Risk: %s\n\n", riskLabel(result.Risk)); err != nil {
		return err
	}
	if opts.Duration > 0 {
		if _, err := fmt.Fprintf(w, "Scan completed in %s\n\n", formatDuration(opts.Duration)); err != nil {
			return err
		}
	}

	for _, check := range result.Checks {
		line := fmt.Sprintf("%s %s", marker(check.Status, opts.NoColor), check.Title)
		if opts.Verbose {
			line = fmt.Sprintf("%s %d/%d", line, check.Points, check.MaxPoints)
		}
		if _, err := fmt.Fprintf(w, "%s\n  %s\n", line, check.Summary); err != nil {
			return err
		}
	}

	weaknesses := topWeaknesses(result.Checks)
	if len(weaknesses) > 0 {
		if _, err := fmt.Fprintln(w, "\nTop weaknesses:"); err != nil {
			return err
		}
		for _, weakness := range weaknesses {
			if _, err := fmt.Fprintf(w, "• %s\n", weakness); err != nil {
				return err
			}
		}
	}

	if len(result.Recommendations) > 0 {
		if _, err := fmt.Fprintln(w, "\nTop fixes:"); err != nil {
			return err
		}
		for i, rec := range result.Recommendations {
			if _, err := fmt.Fprintf(w, "%d. %s - %s\n", i+1, rec.Title, rec.Detail); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeCompact(w io.Writer, result model.ScanResult) error {
	if _, err := fmt.Fprintf(w, "Score: %d/%d (%s)\n", result.Score, result.MaxScore, result.Grade); err != nil {
		return err
	}
	issues := compactIssues(result.Checks)
	if len(issues) == 0 {
		_, err := fmt.Fprintln(w, "Issues: none")
		return err
	}
	_, err := fmt.Fprintf(w, "Issues: %s\n", strings.Join(issues, ", "))
	return err
}

func marker(status model.Status, noColor bool) string {
	text := strings.ToUpper(string(status))
	if noColor {
		return text
	}
	switch status {
	case model.StatusPass:
		return "\x1b[32m✓ PASS\x1b[0m"
	case model.StatusWarn:
		return "\x1b[33m⚠ WARN\x1b[0m"
	case model.StatusFail:
		return "\x1b[31m✗ FAIL\x1b[0m"
	default:
		return "\x1b[36m• INFO\x1b[0m"
	}
}

func riskLabel(risk model.RiskLevel) string {
	switch risk {
	case model.RiskLow:
		return "Low"
	case model.RiskMedium:
		return "Medium"
	case model.RiskHigh:
		return "High"
	default:
		return strings.ToUpper(string(risk))
	}
}

func formatDuration(duration time.Duration) string {
	if duration < time.Second {
		return duration.Round(time.Millisecond).String()
	}
	return duration.Round(100 * time.Millisecond).String()
}

func topWeaknesses(checks []model.CheckResult) []string {
	problems := problemChecks(checks)
	if len(problems) > 3 {
		problems = problems[:3]
	}
	out := make([]string, 0, len(problems))
	for _, check := range problems {
		out = append(out, weaknessLabel(check))
	}
	return out
}

func compactIssues(checks []model.CheckResult) []string {
	problems := problemChecks(checks)
	out := make([]string, 0, len(problems))
	for _, check := range problems {
		out = append(out, compactIssueLabel(check))
	}
	return out
}

func problemChecks(checks []model.CheckResult) []model.CheckResult {
	out := make([]model.CheckResult, 0)
	for _, check := range checks {
		if check.Status.IsProblem() {
			out = append(out, check)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		leftMissing := out[i].MaxPoints - out[i].Points
		rightMissing := out[j].MaxPoints - out[j].Points
		if leftMissing != rightMissing {
			return leftMissing > rightMissing
		}
		return statusRank(out[i].Status) > statusRank(out[j].Status)
	})
	return out
}

func statusRank(status model.Status) int {
	switch status {
	case model.StatusFail:
		return 3
	case model.StatusWarn:
		return 2
	default:
		return 1
	}
}

func weaknessLabel(check model.CheckResult) string {
	switch check.ID {
	case "repository.archived":
		return "Repository archived"
	case "readme.quality":
		return "README needs work"
	case "commits.activity":
		return "Low commit activity"
	case "issues.health":
		return "Stale issues"
	case "prs.health":
		return "Stale pull requests"
	case "ci.present":
		return "CI missing"
	case "license.present":
		return "No license"
	case "releases.present":
		return "No releases"
	case "tests.hints":
		return "Missing tests"
	case "dependencies.hints":
		return "Dependency hygiene"
	default:
		return check.Title
	}
}

func compactIssueLabel(check model.CheckResult) string {
	switch check.ID {
	case "repository.archived":
		return "repo archived"
	case "readme.quality":
		return "README thin"
	case "commits.activity":
		return "low commit activity"
	case "issues.health":
		return "stale issues"
	case "prs.health":
		return "stale PRs"
	case "ci.present":
		return "CI missing"
	case "license.present":
		return "no license"
	case "releases.present":
		return "no releases"
	case "tests.hints":
		return "missing tests"
	case "dependencies.hints":
		return "dependency hygiene"
	default:
		return strings.ToLower(check.Title)
	}
}
