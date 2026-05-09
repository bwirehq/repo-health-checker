package scoring

import (
	"sort"

	"github.com/bwirehq/repo-health-checker/internal/model"
)

const archivedPenalty = 10

func Aggregate(data model.RepositoryData, checks []model.CheckResult, maxRecommendations int) model.ScanResult {
	total, max := 0, 0
	for _, check := range checks {
		total += check.Points
		max += check.MaxPoints
	}

	score := 0
	if max > 0 {
		score = int(float64(total)/float64(max)*100 + 0.5)
	}
	if data.Archived {
		score -= archivedPenalty
		if score < 0 {
			score = 0
		}
	}
	risk := Risk(score, data.Archived, checks)

	return model.ScanResult{
		Repository:      repositorySummary(data),
		Score:           score,
		MaxScore:        100,
		Grade:           Grade(score),
		Risk:            risk,
		Checks:          checks,
		Recommendations: recommendations(checks, maxRecommendations),
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

func Risk(score int, archived bool, checks []model.CheckResult) model.RiskLevel {
	if archived || score < 60 || hasFailedCriticalCheck(checks) {
		return model.RiskHigh
	}
	if score < 80 || hasProblemCheck(checks) {
		return model.RiskMedium
	}
	return model.RiskLow
}

func hasFailedCriticalCheck(checks []model.CheckResult) bool {
	for _, check := range checks {
		if check.Status == model.StatusFail && (check.ID == "ci.present" || check.ID == "commits.activity") {
			return true
		}
	}
	return false
}

func hasProblemCheck(checks []model.CheckResult) bool {
	for _, check := range checks {
		if check.Status.IsProblem() {
			return true
		}
	}
	return false
}

func repositorySummary(data model.RepositoryData) model.RepositorySummary {
	return model.RepositorySummary{
		Owner:         data.Ref.Owner,
		Name:          data.Ref.Name,
		Description:   data.Description,
		DefaultBranch: data.DefaultBranch,
		Archived:      data.Archived,
		Stars:         data.StarCount,
		Forks:         data.ForkCount,
		OpenIssues:    data.OpenIssueCount,
	}
}

func recommendations(checks []model.CheckResult, limit int) []model.Recommendation {
	if limit <= 0 {
		return nil
	}
	candidates := make([]model.CheckResult, 0)
	for _, check := range checks {
		if check.Recommendation != nil {
			candidates = append(candidates, check)
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		leftMissing := candidates[i].MaxPoints - candidates[i].Points
		rightMissing := candidates[j].MaxPoints - candidates[j].Points
		if leftMissing != rightMissing {
			return leftMissing > rightMissing
		}
		return statusRank(candidates[i].Status) > statusRank(candidates[j].Status)
	})

	out := make([]model.Recommendation, 0, min(limit, len(candidates)))
	for _, check := range candidates {
		if len(out) == limit {
			break
		}
		out = append(out, *check.Recommendation)
	}
	return out
}

func statusRank(status model.Status) int {
	switch status {
	case model.StatusFail:
		return 3
	case model.StatusWarn:
		return 2
	case model.StatusInfo:
		return 1
	default:
		return 0
	}
}
