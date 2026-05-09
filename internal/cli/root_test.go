package cli

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestExecuteRejectsInvalidRepo(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute(context.Background(), []string{"scan", "https://example.com/not/repo"}, strings.NewReader(""), &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "only github.com") {
		t.Fatalf("stderr did not explain invalid host: %s", stderr.String())
	}
}

func TestRepoInputPromptsWhenMissingArgument(t *testing.T) {
	var stdout bytes.Buffer
	got, err := repoInput(strings.NewReader("https://github.com/github/cli\n"), &stdout, nil)
	if err != nil {
		t.Fatalf("repoInput returned error: %v", err)
	}
	if got != "https://github.com/github/cli" {
		t.Fatalf("input = %q, want GitHub URL", got)
	}
	if !strings.Contains(stdout.String(), "GitHub repository") {
		t.Fatalf("prompt was not written: %q", stdout.String())
	}
}

func TestScanCommandHasCompactFlag(t *testing.T) {
	cmd := newRootCommand(context.Background(), strings.NewReader(""), io.Discard, io.Discard)
	var scan *cobra.Command
	for _, child := range cmd.Commands() {
		if child.Name() == "scan" {
			scan = child
			break
		}
	}
	if scan == nil {
		t.Fatal("scan command was not registered")
	}
	if flag := scan.Flags().Lookup("compact"); flag == nil {
		t.Fatal("compact flag was not registered")
	}
}
