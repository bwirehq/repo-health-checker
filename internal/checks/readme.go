package checks

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
)

type ReadmeCheck struct{}

func (ReadmeCheck) Run(_ context.Context, data model.RepositoryData, _ config.Config) model.CheckResult {
	const id = "readme.quality"
	readme := strings.TrimSpace(data.Readme)
	if readme == "" {
		return result(id, "README quality", model.StatusFail, 0, 15, "No README was found.", recommendation(id, "Add a README", "Document what the project does, how to install it, and how to run tests."))
	}

	score := 5
	signals := 0
	lower := strings.ToLower(readme)
	if utf8.RuneCountInString(readme) >= 800 {
		score += 4
		signals++
	}
	if strings.Contains(readme, "#") {
		score += 2
		signals++
	}
	for _, keyword := range []string{"install", "usage", "quick start", "getting started", "test", "contributing"} {
		if strings.Contains(lower, keyword) {
			score += 1
			signals++
		}
	}

	if score >= 13 {
		return result(id, "README quality", model.StatusPass, 15, 15, "README includes enough structure and usage detail to orient contributors.", nil)
	}
	if signals >= 2 {
		return result(id, "README quality", model.StatusWarn, score, 15, "README exists but could explain setup, usage, or contribution workflow more clearly.", recommendation(id, "Improve the README", "Add setup steps, usage examples, test commands, and contribution guidance."))
	}
	return result(id, "README quality", model.StatusWarn, score, 15, "README is present but very thin.", recommendation(id, "Expand the README", "Include project purpose, installation, usage, testing, and support information."))
}
