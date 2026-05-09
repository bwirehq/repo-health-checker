package local

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/bwirehq/repo-health-checker/internal/model"
)

func TestSourceFetchRepositoryFromFilesWithoutGit(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "README.md", "# Project\n\nInstall and usage instructions with tests and contributing guidance.")
	writeFile(t, root, "LICENSE", "MIT")
	writeFile(t, root, "go.mod", "module example.com/project\n")
	writeFile(t, root, "go.sum", "")
	writeFile(t, root, ".github/workflows/ci.yml", "name: ci\n")
	writeFile(t, root, "main_test.go", "package main\n")
	writeFile(t, root, "node_modules/ignored/package.json", "{}")

	data, err := NewSource(root).FetchRepository(context.Background(), RefForPath(root))
	if err != nil {
		t.Fatalf("FetchRepository returned error: %v", err)
	}
	if data.Source != model.SourceLocal {
		t.Fatalf("source = %s, want local", data.Source)
	}
	if data.Ref.Name != filepath.Base(root) {
		t.Fatalf("ref = %#v, want local temp dir name", data.Ref)
	}
	if data.Readme == "" {
		t.Fatal("README was not loaded")
	}
	if len(data.WorkflowFiles) != 1 {
		t.Fatalf("workflow files = %#v, want one", data.WorkflowFiles)
	}
	if len(data.DependencyFiles) != 1 || data.DependencyFiles[0] != "go.mod" {
		t.Fatalf("dependency files = %#v, want go.mod only", data.DependencyFiles)
	}
	if len(data.TestFiles) != 1 || data.TestFiles[0] != "main_test.go" {
		t.Fatalf("test files = %#v, want main_test.go", data.TestFiles)
	}
	for _, file := range data.TreeFiles {
		if file == "node_modules/ignored/package.json" {
			t.Fatal("node_modules was not skipped")
		}
	}
}

func TestSourceUsesGitMetadataWhenAvailable(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not available")
	}
	root := t.TempDir()
	runGit(t, root, "init")
	runGit(t, root, "config", "user.email", "test@example.com")
	runGit(t, root, "config", "user.name", "Test User")
	runGit(t, root, "remote", "add", "origin", "https://github.com/acme/project.git")
	writeFile(t, root, "README.md", "# Project\n")
	runGit(t, root, "add", "README.md")
	runGit(t, root, "commit", "-m", "initial")
	runGit(t, root, "tag", "v1.0.0")

	data, err := NewSource(root).FetchRepository(context.Background(), RefForPath(root))
	if err != nil {
		t.Fatalf("FetchRepository returned error: %v", err)
	}
	if data.Ref != (model.RepoRef{Owner: "acme", Name: "project"}) {
		t.Fatalf("ref = %#v, want acme/project", data.Ref)
	}
	if len(data.Commits) != 1 {
		t.Fatalf("commits = %#v, want one commit", data.Commits)
	}
	if len(data.Tags) != 1 || data.Tags[0].Name != "v1.0.0" {
		t.Fatalf("tags = %#v, want v1.0.0", data.Tags)
	}
	if len(data.Releases) != 1 || data.Releases[0].TagName != "v1.0.0" {
		t.Fatalf("releases = %#v, want tag-backed release", data.Releases)
	}
}

func TestIsPath(t *testing.T) {
	root := t.TempDir()
	if !IsPath(".") || !IsPath(root) || !IsPath("..") {
		t.Fatal("expected local path inputs to be detected")
	}
	if IsPath("github/cli") {
		t.Fatal("owner/repo should not be treated as a local path when it does not exist")
	}
}

func writeFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
}

func runGit(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}
