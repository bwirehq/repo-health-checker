package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bwirehq/repo-health-checker/internal/model"
)

func TestWriteText(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, fixtureResult(), Options{NoColor: true, Verbose: true})
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	got := buf.String()
	for _, want := range []string{"Repo Health: 78/100", "PASS CI configured 10/10", "Top fixes:", "Add a license"} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q:\n%s", want, got)
		}
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, fixtureResult(), Options{JSON: true})
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	got := buf.String()
	for _, want := range []string{`"score": 78`, `"checks":`, `"recommendations":`} {
		if !strings.Contains(got, want) {
			t.Fatalf("json missing %q:\n%s", want, got)
		}
	}
}

func fixtureResult() model.ScanResult {
	rec := model.Recommendation{CheckID: "license.present", Title: "Add a license", Detail: "Choose a license."}
	return model.ScanResult{
		Repository: model.RepoRef{Owner: "openai", Name: "codex"},
		Score:      78,
		MaxScore:   100,
		Grade:      "C",
		Checks: []model.CheckResult{
			{ID: "ci.present", Title: "CI configured", Status: model.StatusPass, Points: 10, MaxPoints: 10, Summary: "Continuous integration configuration was detected."},
			{ID: "license.present", Title: "License", Status: model.StatusFail, Points: 0, MaxPoints: 10, Summary: "No license was detected.", Recommendation: &rec},
		},
		Recommendations: []model.Recommendation{rec},
	}
}
