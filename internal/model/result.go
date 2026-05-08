package model

type ScanResult struct {
	Repository      RepoRef          `json:"repository"`
	Score           int              `json:"score"`
	MaxScore        int              `json:"max_score"`
	Grade           string           `json:"grade"`
	Checks          []CheckResult    `json:"checks"`
	Recommendations []Recommendation `json:"recommendations"`
}

type CheckResult struct {
	ID             string          `json:"id"`
	Title          string          `json:"title"`
	Status         Status          `json:"status"`
	Points         int             `json:"points"`
	MaxPoints      int             `json:"max_points"`
	Summary        string          `json:"summary"`
	Recommendation *Recommendation `json:"recommendation,omitempty"`
}

type Recommendation struct {
	CheckID string `json:"check_id"`
	Title   string `json:"title"`
	Detail  string `json:"detail"`
}
