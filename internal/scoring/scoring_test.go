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
	got := Aggregate(model.RepoRef{Owner: "o", Name: "r"}, checks, 3)
	if got.Score != 75 {
		t.Fatalf("score = %d, want 75", got.Score)
	}
	if got.Grade != "C" {
		t.Fatalf("grade = %s, want C", got.Grade)
	}
	if len(got.Recommendations) != 1 {
		t.Fatalf("recommendations = %d, want 1", len(got.Recommendations))
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
