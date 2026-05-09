package model

type ScanResult struct {
	Repository      RepositorySummary `json:"repository"`
	Score           int               `json:"score"`
	MaxScore        int               `json:"max_score"`
	Grade           string            `json:"grade"`
	Risk            RiskLevel         `json:"risk"`
	Checks          []CheckResult     `json:"checks"`
	Recommendations []Recommendation  `json:"recommendations"`
}

type RepositorySummary struct {
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	DefaultBranch string `json:"default_branch,omitempty"`
	Archived      bool   `json:"archived"`
	Stars         int    `json:"stars"`
	Forks         int    `json:"forks"`
	OpenIssues    int    `json:"open_issues"`
}

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

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
