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
	for _, want := range []string{"Repo Health: 78/100", "Repository Risk: Medium", "PASS CI configured 15/15", "Top fixes:", "Add a license"} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q:\n%s", want, got)
		}
	}
}

func TestWriteTextUsesBadgesByDefault(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, fixtureResult(), Options{})
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	got := buf.String()
	for _, want := range []string{"✓ PASS", "✗ FAIL"} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing badge %q:\n%s", want, got)
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
	for _, want := range []string{`"score": 78`, `"risk": "medium"`, `"repository":`, `"archived": false`, `"checks":`, `"recommendations":`} {
		if !strings.Contains(got, want) {
			t.Fatalf("json missing %q:\n%s", want, got)
		}
	}
}

func fixtureResult() model.ScanResult {
	rec := model.Recommendation{CheckID: "license.present", Title: "Add a license", Detail: "Choose a license."}
	return model.ScanResult{
		Repository: model.RepositorySummary{Owner: "github", Name: "cli", DefaultBranch: "main"},
		Score:      78,
		MaxScore:   100,
		Grade:      "C",
		Risk:       model.RiskMedium,
		Checks: []model.CheckResult{
			{ID: "ci.present", Title: "CI configured", Status: model.StatusPass, Points: 15, MaxPoints: 15, Summary: "Continuous integration configuration was detected."},
			{ID: "license.present", Title: "License", Status: model.StatusFail, Points: 0, MaxPoints: 10, Summary: "No license was detected.", Recommendation: &rec},
		},
		Recommendations: []model.Recommendation{rec},
	}
}
