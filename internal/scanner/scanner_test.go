package scanner

import (
	"context"
	"testing"
	"time"

	"github.com/bwirehq/repo-health-checker/internal/checks"
	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
)

func TestScanUsesSourceAndChecks(t *testing.T) {
	ref := model.RepoRef{Owner: "github", Name: "cli"}
	source := fakeSource{data: model.RepositoryData{Ref: ref, Readme: "# Test"}}
	check := fakeCheck{result: model.CheckResult{
		ID:        "fake.check",
		Title:     "Fake check",
		Status:    model.StatusPass,
		Points:    10,
		MaxPoints: 10,
		Summary:   "ok",
	}}

	scan := New(source, config.Default(time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)), []checks.Check{check})
	got, err := scan.Scan(context.Background(), ref)
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if got.Score != 100 {
		t.Fatalf("score = %d, want 100", got.Score)
	}
	if len(got.Checks) != 1 || got.Checks[0].ID != "fake.check" {
		t.Fatalf("unexpected checks: %#v", got.Checks)
	}
}

type fakeSource struct {
	data model.RepositoryData
}

func (f fakeSource) FetchRepository(context.Context, model.RepoRef) (model.RepositoryData, error) {
	return f.data, nil
}

type fakeCheck struct {
	result model.CheckResult
}

func (f fakeCheck) Run(context.Context, model.RepositoryData, config.Config) model.CheckResult {
	return f.result
}
