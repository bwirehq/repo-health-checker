package config

import "time"

type Config struct {
	Now                time.Time
	CommitWindow       time.Duration
	StaleIssueAge      time.Duration
	StalePullAge       time.Duration
	RecentCommitPass   int
	RecentCommitWarn   int
	MaxRecommendations int
}

func Default(now time.Time) Config {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return Config{
		Now:                now,
		CommitWindow:       90 * 24 * time.Hour,
		StaleIssueAge:      90 * 24 * time.Hour,
		StalePullAge:       30 * 24 * time.Hour,
		RecentCommitPass:   10,
		RecentCommitWarn:   1,
		MaxRecommendations: 3,
	}
}
