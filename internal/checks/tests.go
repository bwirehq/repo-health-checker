package checks

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
)

type TestHintCheck struct{}

func (TestHintCheck) Run(_ context.Context, data model.RepositoryData, _ config.Config) model.CheckResult {
	const id = "tests.hints"
	if len(data.TestFiles) > 0 || hasTestPath(data.TreeFiles) {
		return result(id, "Test coverage hints", model.StatusPass, 10, 10, fmt.Sprintf("%d test-related files or paths were detected.", len(data.TestFiles)), nil)
	}
	return result(id, "Test coverage hints", model.StatusWarn, 2, 10, "No common test files or test directories were detected.", recommendation(id, "Add test coverage signals", "Add tests and document the test command so maintainers and users can verify changes."))
}

func hasTestPath(files []string) bool {
	for _, file := range files {
		normalized := strings.ToLower(file)
		if strings.Contains(normalized, "/test/") || strings.Contains(normalized, "/tests/") || strings.HasSuffix(normalized, "_test.go") || strings.HasSuffix(normalized, ".test.js") || strings.HasSuffix(normalized, ".spec.ts") {
			return true
		}
	}
	return false
}
