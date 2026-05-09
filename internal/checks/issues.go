package checks

import (
	"context"
	"fmt"

	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
)

type IssueHealthCheck struct{}

func (IssueHealthCheck) Run(_ context.Context, data model.RepositoryData, cfg config.Config) model.CheckResult {
	const id = "issues.health"
	open := len(data.Issues)
	if open == 0 {
		return result(id, "Issue health", model.StatusPass, 10, 10, "There are no open issues.", nil)
	}

	cutoff := cfg.Now.Add(-cfg.StaleIssueAge)
	stale := countOlderThan(data.Issues, func(issue model.Issue) bool {
		return !issue.UpdatedAt.IsZero() && issue.UpdatedAt.Before(cutoff)
	})
	ratio := float64(stale) / float64(open)

	switch {
	case ratio <= 0.20:
		return result(id, "Issue health", model.StatusPass, 10, 10, fmt.Sprintf("%d open issues; %d appear stale.", open, stale), nil)
	case ratio <= 0.50:
		return result(id, "Issue health", model.StatusWarn, 5, 10, fmt.Sprintf("%d open issues; %d appear stale.", open, stale), recommendation(id, "Triage stale issues", "Review issues with no updates in 90 days and close, label, or prioritize them."))
	default:
		return result(id, "Issue health", model.StatusFail, 2, 10, fmt.Sprintf("%d open issues; %d appear stale.", open, stale), recommendation(id, "Reduce stale issue backlog", "Close outdated issues and label the rest by priority or status."))
	}
}
