package checks

import (
	"context"
	"testing"
	"time"

	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
)

func TestDefaultChecks(t *testing.T) {
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	data := model.RepositoryData{
		Readme:        "# Project\n\nInstall and usage instructions with tests and contributing guidance. " + longText(),
		TreeFiles:     []string{"LICENSE", ".github/workflows/ci.yml", "go.mod", "go.sum", "main_test.go"},
		WorkflowFiles: []string{".github/workflows/ci.yml"},
		Commits: []model.Commit{
			{SHA: "1", AuthorAt: now.Add(-24 * time.Hour)},
			{SHA: "2", AuthorAt: now.Add(-48 * time.Hour)},
			{SHA: "3", AuthorAt: now.Add(-72 * time.Hour)},
			{SHA: "4", AuthorAt: now.Add(-96 * time.Hour)},
			{SHA: "5", AuthorAt: now.Add(-120 * time.Hour)},
			{SHA: "6", AuthorAt: now.Add(-144 * time.Hour)},
			{SHA: "7", AuthorAt: now.Add(-168 * time.Hour)},
			{SHA: "8", AuthorAt: now.Add(-192 * time.Hour)},
			{SHA: "9", AuthorAt: now.Add(-216 * time.Hour)},
			{SHA: "10", AuthorAt: now.Add(-240 * time.Hour)},
		},
		LicenseSPDX:     "MIT",
		DependencyFiles: []string{"go.mod"},
		TestFiles:       []string{"main_test.go"},
		Releases:        []model.Release{{TagName: "v1.0.0"}},
	}
	cfg := config.Default(now)
	for _, check := range DefaultSuite() {
		got := check.Run(context.Background(), data, cfg)
		if got.ID == "" || got.Title == "" {
			t.Fatalf("check returned incomplete result: %#v", got)
		}
		if got.Points < 0 || got.Points > got.MaxPoints {
			t.Fatalf("%s returned invalid score %d/%d", got.ID, got.Points, got.MaxPoints)
		}
	}
}

func TestIssueHealthStale(t *testing.T) {
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	data := model.RepositoryData{
		Issues: []model.Issue{
			{Number: 1, UpdatedAt: now.Add(-120 * 24 * time.Hour)},
			{Number: 2, UpdatedAt: now.Add(-130 * 24 * time.Hour)},
			{Number: 3, UpdatedAt: now.Add(-2 * 24 * time.Hour)},
		},
	}
	got := IssueHealthCheck{}.Run(context.Background(), data, config.Default(now))
	if got.Status != model.StatusFail {
		t.Fatalf("status = %s, want fail", got.Status)
	}
}

func TestGitHubOnlyChecksAreInfoForLocalScans(t *testing.T) {
	data := model.RepositoryData{Source: model.SourceLocal}
	cfg := config.Default(time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC))

	issue := IssueHealthCheck{}.Run(context.Background(), data, cfg)
	if issue.Status != model.StatusInfo || issue.MaxPoints != 0 {
		t.Fatalf("issue check = %#v, want info 0/0", issue)
	}

	pr := PullRequestHealthCheck{}.Run(context.Background(), data, cfg)
	if pr.Status != model.StatusInfo || pr.MaxPoints != 0 {
		t.Fatalf("pull request check = %#v, want info 0/0", pr)
	}
}

func longText() string {
	out := ""
	for len(out) < 900 {
		out += "This repository has enough documentation detail for maintainers and users. "
	}
	return out
}
