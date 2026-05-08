package scanner

import (
	"context"

	"github.com/bwirehq/repo-health-checker/internal/checks"
	"github.com/bwirehq/repo-health-checker/internal/config"
	"github.com/bwirehq/repo-health-checker/internal/model"
	"github.com/bwirehq/repo-health-checker/internal/scoring"
)

type Source interface {
	FetchRepository(context.Context, model.RepoRef) (model.RepositoryData, error)
}

type Scanner struct {
	source Source
	checks []checks.Check
	cfg    config.Config
}

func New(source Source, cfg config.Config, suite []checks.Check) *Scanner {
	if suite == nil {
		suite = checks.DefaultSuite()
	}
	return &Scanner{source: source, checks: suite, cfg: cfg}
}

func (s *Scanner) Scan(ctx context.Context, ref model.RepoRef) (model.ScanResult, error) {
	data, err := s.source.FetchRepository(ctx, ref)
	if err != nil {
		return model.ScanResult{}, err
	}

	results := make([]model.CheckResult, 0, len(s.checks))
	for _, check := range s.checks {
		results = append(results, check.Run(ctx, data, s.cfg))
	}
	return scoring.Aggregate(ref, results, s.cfg.MaxRecommendations), nil
}
