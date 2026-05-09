package checks

import (
	"context"
	"fmt"

	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
)

type PullRequestHealthCheck struct{}

func (PullRequestHealthCheck) Run(_ context.Context, data model.RepositoryData, cfg config.Config) model.CheckResult {
	const id = "prs.health"
	if data.Source == model.SourceLocal {
		return result(id, "Pull request health", model.StatusInfo, 0, 0, "Pull request health is unavailable for offline local scans.", nil)
	}
	open := len(data.PullRequests)
	if open == 0 {
		return result(id, "Pull request health", model.StatusPass, 10, 10, "There are no open pull requests.", nil)
	}

	cutoff := cfg.Now.Add(-cfg.StalePullAge)
	stale := countOlderThan(data.PullRequests, func(pr model.PullRequest) bool {
		return !pr.UpdatedAt.IsZero() && pr.UpdatedAt.Before(cutoff)
	})
	if stale == 0 {
		return result(id, "Pull request health", model.StatusPass, 10, 10, fmt.Sprintf("%d open pull requests; none appear stale.", open), nil)
	}
	if stale <= 3 {
		return result(id, "Pull request health", model.StatusWarn, 6, 10, fmt.Sprintf("%d open pull requests; %d appear stale.", open, stale), recommendation(id, "Review stale pull requests", "Merge, close, or request updates on pull requests with no activity in 30 days."))
	}
	return result(id, "Pull request health", model.StatusFail, 2, 10, fmt.Sprintf("%d open pull requests; %d appear stale.", open, stale), recommendation(id, "Clear stale pull requests", "Prioritize old pull requests so contributors know whether their work is still wanted."))
}
