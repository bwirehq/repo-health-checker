package github

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwirehq/repo-health-checker/internal/model"
	gogithub "github.com/google/go-github/v64/github"
)

type Client struct {
	api *gogithub.Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	api := gogithub.NewClient(httpClient)
	if token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); token != "" {
		api = api.WithAuthToken(token)
	}
	return &Client{api: api}
}

func (c *Client) FetchRepository(ctx context.Context, ref model.RepoRef) (model.RepositoryData, error) {
	repo, _, err := c.api.Repositories.Get(ctx, ref.Owner, ref.Name)
	if err != nil {
		return model.RepositoryData{}, classify(err, "failed to fetch repository metadata")
	}

	data := model.RepositoryData{
		Ref:            ref,
		Description:    repo.GetDescription(),
		DefaultBranch:  repo.GetDefaultBranch(),
		Archived:       repo.GetArchived(),
		OpenIssueCount: repo.GetOpenIssuesCount(),
		StarCount:      repo.GetStargazersCount(),
		ForkCount:      repo.GetForksCount(),
	}
	if repo.CreatedAt != nil {
		data.CreatedAt = repo.CreatedAt.Time
	}
	if repo.PushedAt != nil {
		data.PushedAt = repo.PushedAt.Time
	}
	if repo.License != nil {
		data.LicenseSPDX = repo.License.GetSPDXID()
		data.LicenseName = repo.License.GetName()
	}

	data.Readme = c.fetchReadme(ctx, ref)
	data.TreeFiles = c.fetchTreeFiles(ctx, ref, data.DefaultBranch)
	data.WorkflowFiles = workflowFiles(data.TreeFiles)
	data.DependencyFiles = dependencyFiles(data.TreeFiles)
	data.TestFiles = testFiles(data.TreeFiles)
	data.Commits = c.fetchCommits(ctx, ref)
	data.Issues = c.fetchIssues(ctx, ref)
	data.PullRequests = c.fetchPullRequests(ctx, ref)
	data.Releases = c.fetchReleases(ctx, ref)
	data.Tags = c.fetchTags(ctx, ref)

	return data, nil
}

func (c *Client) fetchReadme(ctx context.Context, ref model.RepoRef) string {
	readme, _, err := c.api.Repositories.GetReadme(ctx, ref.Owner, ref.Name, nil)
	if err != nil {
		return ""
	}
	content, err := readme.GetContent()
	if err != nil {
		return ""
	}
	return content
}

func (c *Client) fetchTreeFiles(ctx context.Context, ref model.RepoRef, branch string) []string {
	if branch == "" {
		return nil
	}
	tree, _, err := c.api.Git.GetTree(ctx, ref.Owner, ref.Name, branch, true)
	if err != nil || tree == nil {
		return nil
	}
	files := make([]string, 0, len(tree.Entries))
	for _, entry := range tree.Entries {
		if entry.GetType() == "blob" {
			files = append(files, entry.GetPath())
		}
	}
	return files
}

func (c *Client) fetchCommits(ctx context.Context, ref model.RepoRef) []model.Commit {
	commits, _, err := c.api.Repositories.ListCommits(ctx, ref.Owner, ref.Name, &gogithub.CommitsListOptions{
		ListOptions: gogithub.ListOptions{PerPage: 100},
	})
	if err != nil {
		return nil
	}
	out := make([]model.Commit, 0, len(commits))
	for _, commit := range commits {
		item := model.Commit{SHA: commit.GetSHA()}
		if commit.Commit != nil && commit.Commit.Author != nil && commit.Commit.Author.Date != nil {
			item.AuthorAt = commit.Commit.Author.Date.Time
		}
		out = append(out, item)
	}
	return out
}

func (c *Client) fetchIssues(ctx context.Context, ref model.RepoRef) []model.Issue {
	issues, _, err := c.api.Issues.ListByRepo(ctx, ref.Owner, ref.Name, &gogithub.IssueListByRepoOptions{
		State:       "open",
		ListOptions: gogithub.ListOptions{PerPage: 100},
	})
	if err != nil {
		return nil
	}
	out := make([]model.Issue, 0, len(issues))
	for _, issue := range issues {
		if issue.IsPullRequest() {
			continue
		}
		out = append(out, model.Issue{
			Number:    issue.GetNumber(),
			Title:     issue.GetTitle(),
			CreatedAt: githubTime(issue.CreatedAt),
			UpdatedAt: githubTime(issue.UpdatedAt),
		})
	}
	return out
}

func (c *Client) fetchPullRequests(ctx context.Context, ref model.RepoRef) []model.PullRequest {
	pulls, _, err := c.api.PullRequests.List(ctx, ref.Owner, ref.Name, &gogithub.PullRequestListOptions{
		State:       "open",
		ListOptions: gogithub.ListOptions{PerPage: 100},
	})
	if err != nil {
		return nil
	}
	out := make([]model.PullRequest, 0, len(pulls))
	for _, pull := range pulls {
		out = append(out, model.PullRequest{
			Number:    pull.GetNumber(),
			Title:     pull.GetTitle(),
			CreatedAt: githubTime(pull.CreatedAt),
			UpdatedAt: githubTime(pull.UpdatedAt),
		})
	}
	return out
}

func (c *Client) fetchReleases(ctx context.Context, ref model.RepoRef) []model.Release {
	releases, _, err := c.api.Repositories.ListReleases(ctx, ref.Owner, ref.Name, &gogithub.ListOptions{PerPage: 20})
	if err != nil {
		return nil
	}
	out := make([]model.Release, 0, len(releases))
	for _, release := range releases {
		out = append(out, model.Release{
			Name:      release.GetName(),
			TagName:   release.GetTagName(),
			CreatedAt: githubTime(release.CreatedAt),
		})
	}
	return out
}

func (c *Client) fetchTags(ctx context.Context, ref model.RepoRef) []model.Tag {
	tags, _, err := c.api.Repositories.ListTags(ctx, ref.Owner, ref.Name, &gogithub.ListOptions{PerPage: 20})
	if err != nil {
		return nil
	}
	out := make([]model.Tag, 0, len(tags))
	for _, tag := range tags {
		out = append(out, model.Tag{Name: tag.GetName()})
	}
	return out
}

func githubTime(value *gogithub.Timestamp) time.Time {
	if value == nil {
		return time.Time{}
	}
	return value.Time
}

func classify(err error, message string) error {
	var rateLimitErr *gogithub.RateLimitError
	if errors.As(err, &rateLimitErr) {
		return wrap(ErrRateLimited, message+": GitHub API rate limit exceeded", err)
	}
	var abuseErr *gogithub.AbuseRateLimitError
	if errors.As(err, &abuseErr) {
		return wrap(ErrRateLimited, message+": GitHub API abuse limit exceeded", err)
	}
	var acceptedErr *gogithub.AcceptedError
	if errors.As(err, &acceptedErr) {
		return wrap(ErrUnavailable, message+": GitHub is still processing this repository", err)
	}
	var responseErr *gogithub.ErrorResponse
	if errors.As(err, &responseErr) {
		switch responseErr.Response.StatusCode {
		case http.StatusNotFound:
			return wrap(ErrNotFound, message+": repository was not found", err)
		case http.StatusForbidden, http.StatusUnauthorized:
			return wrap(ErrForbidden, message+": GitHub denied access", err)
		}
	}
	return wrap(ErrUnavailable, message, err)
}
