package scoring

import "github.com/bwirehq/repo-health-checker/internal/model"

func Aggregate(repo model.RepoRef, checks []model.CheckResult, maxRecommendations int) model.ScanResult {
	total, max := 0, 0
	recommendations := make([]model.Recommendation, 0, maxRecommendations)
	for _, check := range checks {
		total += check.Points
		max += check.MaxPoints
		if check.Recommendation != nil && len(recommendations) < maxRecommendations {
			recommendations = append(recommendations, *check.Recommendation)
		}
	}

	score := 0
	if max > 0 {
		score = int(float64(total)/float64(max)*100 + 0.5)
	}

	return model.ScanResult{
		Repository:      repo,
		Score:           score,
		MaxScore:        100,
		Grade:           Grade(score),
		Checks:          checks,
		Recommendations: recommendations,
	}
}

func Grade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}
