package checks

import (
	"context"
	"fmt"

	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
)

type DependencyHintCheck struct{}

func (DependencyHintCheck) Run(_ context.Context, data model.RepositoryData, _ config.Config) model.CheckResult {
	const id = "dependencies.hints"
	if len(data.DependencyFiles) == 0 {
		return result(id, "Dependency hygiene", model.StatusInfo, 5, 5, "No dependency manifests were detected.", nil)
	}
	if hasLockfile(data.TreeFiles) {
		return result(id, "Dependency hygiene", model.StatusPass, 5, 5, fmt.Sprintf("%d dependency manifests and a lockfile were detected.", len(data.DependencyFiles)), nil)
	}
	return result(id, "Dependency hygiene", model.StatusWarn, 2, 5, fmt.Sprintf("%d dependency manifests were detected, but no common lockfile.", len(data.DependencyFiles)), recommendation(id, "Add or commit lockfiles", "Commit the appropriate lockfile for reproducible installs when the ecosystem expects one."))
}

func hasLockfile(files []string) bool {
	return containsPath(files, "package-lock.json", "pnpm-lock.yaml", "yarn.lock", "Cargo.lock", "go.sum", "poetry.lock", "Pipfile.lock", "Gemfile.lock", "composer.lock")
}
