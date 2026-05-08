package checks

import (
	"context"

	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
)

type LicenseCheck struct{}

func (LicenseCheck) Run(_ context.Context, data model.RepositoryData, _ config.Config) model.CheckResult {
	const id = "license.present"
	if data.LicenseSPDX != "" || containsPath(data.TreeFiles, "LICENSE", "LICENSE.md", "COPYING") {
		return result(id, "License", model.StatusPass, 10, 10, "A license was detected.", nil)
	}
	return result(id, "License", model.StatusFail, 0, 10, "No license file or SPDX license metadata was detected.", recommendation(id, "Add a license", "Choose and add a license so users know how they can use the project."))
}
