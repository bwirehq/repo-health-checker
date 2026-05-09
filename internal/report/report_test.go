package report

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/bwirehq/repo-health-checker/internal/model"
)

func TestWriteText(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, fixtureResult(), Options{NoColor: true, Verbose: true})
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	got := buf.String()
	for _, want := range []string{"Repo Health: 78/100", "Repository Risk: Medium", "Score breakdown:", "License: 0/10", "PASS CI configured 15/15", "Top weaknesses:", "No license", "Top fixes:", "Add a license"} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q:\n%s", want, got)
		}
	}
}

func TestWriteTextIncludesDuration(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, fixtureResult(), Options{NoColor: true, Duration: 1200 * time.Millisecond})
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "Scan completed in 1.2s") {
		t.Fatalf("output missing duration:\n%s", got)
	}
}

func TestWriteTextUsesBadgesByDefault(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, fixtureResult(), Options{})
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	got := buf.String()
	for _, want := range []string{"\u2713 PASS", "\u2717 FAIL"} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing badge %q:\n%s", want, got)
		}
	}
}

func TestWriteCompact(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, fixtureResult(), Options{Compact: true, Duration: 1200 * time.Millisecond})
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	got := buf.String()
	for _, want := range []string{"Score: 78/100 (C)", "Issues: no license, no releases, missing tests"} {
		if !strings.Contains(got, want) {
			t.Fatalf("compact output missing %q:\n%s", want, got)
		}
	}
	for _, notWant := range []string{"Repository Risk:", "Scan completed in", "Score breakdown:", "Top fixes:"} {
		if strings.Contains(got, notWant) {
			t.Fatalf("compact output included %q:\n%s", notWant, got)
		}
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, fixtureResult(), Options{JSON: true, Compact: true, Duration: 1200 * time.Millisecond})
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	got := buf.String()
	for _, want := range []string{`"score": 78`, `"risk": "medium"`, `"repository":`, `"archived": false`, `"checks":`, `"recommendations":`} {
		if !strings.Contains(got, want) {
			t.Fatalf("json missing %q:\n%s", want, got)
		}
	}
	for _, notWant := range []string{"Scan completed in", "Score: 78/100", "Score breakdown:"} {
		if strings.Contains(got, notWant) {
			t.Fatalf("json output included text renderer content %q:\n%s", notWant, got)
		}
	}
}

func fixtureResult() model.ScanResult {
	rec := model.Recommendation{CheckID: "license.present", Title: "Add a license", Detail: "Choose a license."}
	releaseRec := model.Recommendation{CheckID: "releases.present", Title: "Create a release", Detail: "Tag a version."}
	testRec := model.Recommendation{CheckID: "tests.hints", Title: "Add tests", Detail: "Add common test files."}
	return model.ScanResult{
		Repository: model.RepositorySummary{Owner: "github", Name: "cli", DefaultBranch: "main"},
		Score:      78,
		MaxScore:   100,
		Grade:      "C",
		Risk:       model.RiskMedium,
		Checks: []model.CheckResult{
			{ID: "ci.present", Title: "CI configured", Status: model.StatusPass, Points: 15, MaxPoints: 15, Summary: "Continuous integration configuration was detected."},
			{ID: "license.present", Title: "License", Status: model.StatusFail, Points: 0, MaxPoints: 10, Summary: "No license was detected.", Recommendation: &rec},
			{ID: "releases.present", Title: "Releases", Status: model.StatusFail, Points: 0, MaxPoints: 10, Summary: "No releases or tags were found.", Recommendation: &releaseRec},
			{ID: "tests.hints", Title: "Test coverage hints", Status: model.StatusWarn, Points: 2, MaxPoints: 10, Summary: "No common test files or test directories were detected.", Recommendation: &testRec},
		},
		Recommendations: []model.Recommendation{rec, releaseRec, testRec},
	}
}
