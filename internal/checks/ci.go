package checks

import (
	"context"

	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
)

type CICheck struct{}

func (CICheck) Run(_ context.Context, data model.RepositoryData, _ config.Config) model.CheckResult {
	const id = "ci.present"
	if len(data.WorkflowFiles) > 0 || containsPath(data.TreeFiles, ".travis.yml", "circle.yml", "azure-pipelines.yml", "Jenkinsfile") || hasPrefix(data.TreeFiles, ".circleci/", ".buildkite/") {
		return result(id, "CI configured", model.StatusPass, 10, 10, "Continuous integration configuration was detected.", nil)
	}
	return result(id, "CI configured", model.StatusFail, 0, 10, "No common CI configuration was detected.", recommendation(id, "Add CI", "Run tests and lint checks on every pull request with GitHub Actions or another CI provider."))
}
