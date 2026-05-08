package checks

import (
	"context"
	"fmt"

	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
)

type CommitActivityCheck struct{}

func (CommitActivityCheck) Run(_ context.Context, data model.RepositoryData, cfg config.Config) model.CheckResult {
	const id = "commits.activity"
	cutoff := cfg.Now.Add(-cfg.CommitWindow)
	recent := 0
	for _, commit := range data.Commits {
		if !commit.AuthorAt.IsZero() && commit.AuthorAt.After(cutoff) {
			recent++
		}
	}

	if recent >= cfg.RecentCommitPass {
		return result(id, "Commit activity", model.StatusPass, 15, 15, fmt.Sprintf("%d commits were found in the last 90 days.", recent), nil)
	}
	if recent >= cfg.RecentCommitWarn {
		return result(id, "Commit activity", model.StatusWarn, 8, 15, fmt.Sprintf("%d commits were found in the last 90 days.", recent), recommendation(id, "Increase maintenance cadence", "Keep regular commits or maintenance notes so users can see the project is active."))
	}
	return result(id, "Commit activity", model.StatusFail, 0, 15, "No commits were found in the last 90 days.", recommendation(id, "Show recent maintenance", "Merge updates, close maintenance tasks, or archive the repository if it is no longer maintained."))
}
