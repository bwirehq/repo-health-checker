package checks

import (
	"context"

	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
)

type Check interface {
	Run(context.Context, model.RepositoryData, config.Config) model.CheckResult
}

func result(id, title string, status model.Status, points, max int, summary string, rec *model.Recommendation) model.CheckResult {
	return model.CheckResult{
		ID:             id,
		Title:          title,
		Status:         status,
		Points:         clamp(points, 0, max),
		MaxPoints:      max,
		Summary:        summary,
		Recommendation: rec,
	}
}

func recommendation(checkID, title, detail string) *model.Recommendation {
	return &model.Recommendation{CheckID: checkID, Title: title, Detail: detail}
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
