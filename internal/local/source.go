package local

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/bwirehq/repo-health-checker/internal/model"
)

type Source struct {
	root string
}

func NewSource(root string) *Source {
	return &Source{root: root}
}

func IsPath(input string) bool {
	value := strings.TrimSpace(input)
	if value == "." || strings.HasPrefix(value, "..") || filepath.IsAbs(value) || strings.Contains(value, `\`) {
		return true
	}
	if _, err := os.Stat(value); err == nil {
		return true
	}
	return false
}

func RefForPath(path string) model.RepoRef {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	name := filepath.Base(abs)
	if name == "." || name == string(filepath.Separator) || name == "" {
		name = "local"
	}
	return model.RepoRef{Owner: "local", Name: sanitizePart(name)}
}

func (s *Source) FetchRepository(ctx context.Context, ref model.RepoRef) (model.RepositoryData, error) {
	root, err := filepath.Abs(s.root)
	if err != nil {
		return model.RepositoryData{}, err
	}
	stat, err := os.Stat(root)
	if err != nil {
		return model.RepositoryData{}, err
	}
	if !stat.IsDir() {
		root = filepath.Dir(root)
	}

	files := walkFiles(root)
	readme := readFirst(root, "README.md", "README", "README.txt", "README.rst")
	tags := gitTags(ctx, root)

	data := model.RepositoryData{
		Ref:             remoteRef(ctx, root, ref),
		Source:          model.SourceLocal,
		DefaultBranch:   gitOutput(ctx, root, "branch", "--show-current"),
		Archived:        false,
		Readme:          readme,
		TreeFiles:       files,
		WorkflowFiles:   workflowFiles(files),
		DependencyFiles: dependencyFiles(files),
		TestFiles:       testFiles(files),
		Commits:         gitCommits(ctx, root),
		Tags:            tags,
		Releases:        releasesFromTags(tags),
	}
	return data, nil
}

func walkFiles(root string) []string {
	var files []string
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.IsDir() {
			if shouldSkipDir(entry.Name()) && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	return files
}

func shouldSkipDir(name string) bool {
	switch strings.ToLower(name) {
	case ".git", "bin", "node_modules", "vendor", "dist", "build", ".cache", "coverage":
		return true
	default:
		return false
	}
}

func readFirst(root string, names ...string) string {
	for _, name := range names {
		content, err := os.ReadFile(filepath.Join(root, name))
		if err == nil {
			return string(content)
		}
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return ""
	}
	candidates := make(map[string]struct{}, len(names))
	for _, name := range names {
		candidates[strings.ToLower(name)] = struct{}{}
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if _, ok := candidates[strings.ToLower(entry.Name())]; !ok {
			continue
		}
		content, err := os.ReadFile(filepath.Join(root, entry.Name()))
		if err == nil {
			return string(content)
		}
	}
	return ""
}

func gitCommits(ctx context.Context, root string) []model.Commit {
	out := gitOutput(ctx, root, "log", "--since=90.days", "--format=%H%x09%cI")
	if out == "" {
		return nil
	}
	var commits []model.Commit
	for _, line := range strings.Split(out, "\n") {
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			continue
		}
		authorAt, err := time.Parse(time.RFC3339, parts[1])
		if err != nil {
			continue
		}
		commits = append(commits, model.Commit{SHA: parts[0], AuthorAt: authorAt})
	}
	return commits
}

func gitTags(ctx context.Context, root string) []model.Tag {
	out := gitOutput(ctx, root, "tag", "--list")
	if out == "" {
		return nil
	}
	var tags []model.Tag
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			tags = append(tags, model.Tag{Name: line})
		}
	}
	return tags
}

func releasesFromTags(tags []model.Tag) []model.Release {
	releases := make([]model.Release, 0, len(tags))
	for _, tag := range tags {
		releases = append(releases, model.Release{TagName: tag.Name})
	}
	return releases
}

func remoteRef(ctx context.Context, root string, fallback model.RepoRef) model.RepoRef {
	remote := gitOutput(ctx, root, "config", "--get", "remote.origin.url")
	if remote == "" {
		return fallback
	}
	remote = strings.TrimSuffix(strings.TrimSpace(remote), ".git")
	remote = strings.TrimPrefix(remote, "git@github.com:")
	remote = strings.TrimPrefix(remote, "https://github.com/")
	parts := strings.Split(remote, "/")
	if len(parts) < 2 {
		return fallback
	}
	owner, name := sanitizePart(parts[len(parts)-2]), sanitizePart(parts[len(parts)-1])
	if owner == "" || name == "" {
		return fallback
	}
	return model.RepoRef{Owner: owner, Name: name}
}

func gitOutput(ctx context.Context, root string, args ...string) string {
	fullArgs := append([]string{"-C", root}, args...)
	out, err := exec.CommandContext(ctx, "git", fullArgs...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func sanitizePart(value string) string {
	value = strings.TrimSpace(value)
	var b strings.Builder
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '.', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteRune('-')
		}
	}
	return strings.Trim(b.String(), ".-")
}

func workflowFiles(files []string) []string {
	var out []string
	for _, file := range files {
		normalized := strings.ToLower(file)
		if strings.HasPrefix(normalized, ".github/workflows/") && (strings.HasSuffix(normalized, ".yml") || strings.HasSuffix(normalized, ".yaml")) {
			out = append(out, file)
		}
	}
	return out
}

func dependencyFiles(files []string) []string {
	names := map[string]struct{}{
		"package.json":     {},
		"go.mod":           {},
		"pyproject.toml":   {},
		"requirements.txt": {},
		"cargo.toml":       {},
		"gemfile":          {},
		"composer.json":    {},
		"pom.xml":          {},
		"build.gradle":     {},
	}
	return filterByName(files, names)
}

func testFiles(files []string) []string {
	var out []string
	for _, file := range files {
		normalized := strings.ToLower(file)
		if strings.Contains(normalized, "/test/") ||
			strings.Contains(normalized, "/tests/") ||
			strings.HasSuffix(normalized, "_test.go") ||
			strings.HasSuffix(normalized, ".test.js") ||
			strings.HasSuffix(normalized, ".spec.js") ||
			strings.HasSuffix(normalized, ".test.ts") ||
			strings.HasSuffix(normalized, ".spec.ts") ||
			strings.HasSuffix(normalized, ".test.tsx") ||
			strings.HasSuffix(normalized, ".spec.tsx") {
			out = append(out, file)
		}
	}
	return out
}

func filterByName(files []string, names map[string]struct{}) []string {
	var out []string
	for _, file := range files {
		parts := strings.Split(strings.ToLower(file), "/")
		name := parts[len(parts)-1]
		if _, ok := names[name]; ok {
			out = append(out, file)
		}
	}
	return out
}
