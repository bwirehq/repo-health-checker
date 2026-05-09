package model

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var repoPartPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

type RepoRef struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

func (r RepoRef) String() string {
	return r.Owner + "/" + r.Name
}

func ParseRepoRef(input string) (RepoRef, error) {
	value := strings.TrimSpace(input)
	if value == "" {
		return RepoRef{}, errors.New("repository is required")
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		parsed, err := url.Parse(value)
		if err != nil {
			return RepoRef{}, fmt.Errorf("invalid GitHub URL: %w", err)
		}
		if parsed.Scheme != "https" {
			return RepoRef{}, errors.New("only https://github.com URLs are supported")
		}
		if !strings.EqualFold(parsed.Host, "github.com") {
			return RepoRef{}, errors.New("only github.com repository URLs are supported")
		}
		value = strings.Trim(parsed.Path, "/")
	}

	parts := strings.Split(value, "/")
	if len(parts) != 2 {
		return RepoRef{}, errors.New("repository must be in owner/repo format")
	}
	owner, name := parts[0], strings.TrimSuffix(parts[1], ".git")
	if !validRepoPart(owner) || !validRepoPart(name) {
		return RepoRef{}, errors.New("repository owner and name may only contain letters, numbers, dot, dash, and underscore")
	}
	return RepoRef{Owner: owner, Name: name}, nil
}

func validRepoPart(value string) bool {
	return value != "" && repoPartPattern.MatchString(value) && value != "." && value != ".."
}

type RepositoryData struct {
	Ref             RepoRef
	Source          SourceType
	Description     string
	DefaultBranch   string
	Archived        bool
	CreatedAt       time.Time
	PushedAt        time.Time
	OpenIssueCount  int
	StarCount       int
	ForkCount       int
	Readme          string
	LicenseSPDX     string
	LicenseName     string
	TreeFiles       []string
	WorkflowFiles   []string
	Commits         []Commit
	Issues          []Issue
	PullRequests    []PullRequest
	Releases        []Release
	Tags            []Tag
	DependencyFiles []string
	TestFiles       []string
}

type SourceType string

const (
	SourceGitHub SourceType = "github"
	SourceLocal  SourceType = "local"
)

type Commit struct {
	SHA      string
	AuthorAt time.Time
}

type Issue struct {
	Number    int
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PullRequest struct {
	Number    int
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Release struct {
	Name      string
	TagName   string
	CreatedAt time.Time
}

type Tag struct {
	Name string
}
