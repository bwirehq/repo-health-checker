package checks

import (
	"context"

	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
)

type ArchivedCheck struct{}

func (ArchivedCheck) Run(_ context.Context, data model.RepositoryData, _ config.Config) model.CheckResult {
	const id = "repository.archived"
	if !data.Archived {
		return result(id, "Repository status", model.StatusPass, 0, 0, "Repository is active.", nil)
	}
	return result(id, "Repository status", model.StatusWarn, 0, 0, "Repository is archived and read-only.", recommendation(id, "Clarify maintenance status", "Unarchive the repository if work continues, or document that users should treat it as unmaintained."))
}
