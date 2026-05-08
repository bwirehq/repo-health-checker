package checks

import (
	"context"
	"fmt"

	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
)

type ReleaseCheck struct{}

func (ReleaseCheck) Run(_ context.Context, data model.RepositoryData, _ config.Config) model.CheckResult {
	const id = "releases.present"
	if len(data.Releases) > 0 {
		return result(id, "Releases", model.StatusPass, 10, 10, fmt.Sprintf("%d releases were found.", len(data.Releases)), nil)
	}
	if len(data.Tags) > 0 {
		return result(id, "Releases", model.StatusWarn, 6, 10, fmt.Sprintf("%d tags were found, but no GitHub releases.", len(data.Tags)), recommendation(id, "Create GitHub releases", "Publish release notes from existing tags so users can understand changes."))
	}
	return result(id, "Releases", model.StatusFail, 0, 10, "No releases or tags were found.", recommendation(id, "Create a release", "Tag a version and publish release notes when the project reaches a usable milestone."))
}
