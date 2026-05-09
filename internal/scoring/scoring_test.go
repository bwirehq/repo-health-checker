package scoring

import (
	"testing"

	"github.com/bwirehq/repo-health-checker/internal/model"
)

func TestAggregate(t *testing.T) {
	checks := []model.CheckResult{
		{ID: "a", Points: 10, MaxPoints: 10},
		{ID: "b", Points: 5, MaxPoints: 10, Recommendation: &model.Recommendation{CheckID: "b", Title: "Fix b"}},
	}
	got := Aggregate(model.RepositoryData{Ref: model.RepoRef{Owner: "o", Name: "r"}}, checks, 3)
	if got.Score != 75 {
		t.Fatalf("score = %d, want 75", got.Score)
	}
	if got.Grade != "C" {
		t.Fatalf("grade = %s, want C", got.Grade)
	}
	if len(got.Recommendations) != 1 {
		t.Fatalf("recommendations = %d, want 1", len(got.Recommendations))
	}
	if got.Risk != model.RiskMedium {
		t.Fatalf("risk = %s, want medium", got.Risk)
	}
	if got.Repository.Owner != "o" || got.Repository.Name != "r" {
		t.Fatalf("repository = %#v, want o/r", got.Repository)
	}
}

func TestAggregateArchivedRepoAppliesPenaltyAndHighRisk(t *testing.T) {
	checks := []model.CheckResult{
		{ID: "repository.archived", Status: model.StatusWarn, Points: 0, MaxPoints: 0},
		{ID: "ci.present", Status: model.StatusPass, Points: 15, MaxPoints: 15},
		{ID: "commits.activity", Status: model.StatusPass, Points: 20, MaxPoints: 20},
	}
	got := Aggregate(model.RepositoryData{Ref: model.RepoRef{Owner: "o", Name: "r"}, Archived: true}, checks, 3)
	if got.Score != 90 {
		t.Fatalf("score = %d, want 90", got.Score)
	}
	if got.Risk != model.RiskHigh {
		t.Fatalf("risk = %s, want high", got.Risk)
	}
	if !got.Repository.Archived {
		t.Fatal("repository archived flag was not preserved")
	}
}

func TestRecommendationsSortByMissingPoints(t *testing.T) {
	checks := []model.CheckResult{
		{ID: "license.present", Status: model.StatusFail, Points: 0, MaxPoints: 10, Recommendation: &model.Recommendation{CheckID: "license.present", Title: "Add a license"}},
		{ID: "ci.present", Status: model.StatusFail, Points: 0, MaxPoints: 15, Recommendation: &model.Recommendation{CheckID: "ci.present", Title: "Add CI"}},
		{ID: "commits.activity", Status: model.StatusWarn, Points: 10, MaxPoints: 20, Recommendation: &model.Recommendation{CheckID: "commits.activity", Title: "Increase cadence"}},
	}
	got := Aggregate(model.RepositoryData{Ref: model.RepoRef{Owner: "o", Name: "r"}}, checks, 2)
	if len(got.Recommendations) != 2 {
		t.Fatalf("recommendations = %d, want 2", len(got.Recommendations))
	}
	if got.Recommendations[0].CheckID != "ci.present" || got.Recommendations[1].CheckID != "license.present" {
		t.Fatalf("recommendations not sorted by impact: %#v", got.Recommendations)
	}
}

func TestGrade(t *testing.T) {
	tests := map[int]string{95: "A", 85: "B", 75: "C", 65: "D", 20: "F"}
	for score, want := range tests {
		if got := Grade(score); got != want {
			t.Fatalf("Grade(%d) = %s, want %s", score, got, want)
		}
	}
}
