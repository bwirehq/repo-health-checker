package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bwirehq/repo-health-checker/internal/config"
	gh "github.com/bwirehq/repo-health-checker/internal/github"
	"github.com/bwirehq/repo-health-checker/internal/local"
	"github.com/bwirehq/repo-health-checker/internal/model"
	"github.com/bwirehq/repo-health-checker/internal/report"
	"github.com/bwirehq/repo-health-checker/internal/scanner"
	"github.com/spf13/cobra"
)

func Execute(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	root := newRootCommand(ctx, stdin, stdout, stderr)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCode(err)
	}
	return 0
}

func newRootCommand(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	var opts scanOptions
	cmd := &cobra.Command{
		Use:           "repo-health",
		Short:         "Scan GitHub repositories and report a health score.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	scan := &cobra.Command{
		Use:   "scan [owner/repo|github-url]",
		Short: "Scan a public GitHub repository.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := repoInput(stdin, stdout, args)
			if err != nil {
				return err
			}
			return runScan(ctx, stdout, input, opts)
		},
	}
	scan.Flags().BoolVar(&opts.json, "json", false, "write machine-readable JSON")
	scan.Flags().BoolVar(&opts.compact, "compact", false, "write compact text output")
	scan.Flags().BoolVar(&opts.noColor, "no-color", false, "disable ANSI color output")
	scan.Flags().BoolVar(&opts.verbose, "verbose", false, "include score details in text output")
	scan.Flags().IntVar(&opts.failUnder, "fail-under", 0, "exit with code 2 when score is below this value")
	scan.Flags().DurationVar(&opts.timeout, "timeout", 15*time.Second, "GitHub API timeout")
	cmd.AddCommand(scan)
	cmd.SetErr(stderr)
	cmd.SetOut(stdout)
	return cmd
}

func repoInput(stdin io.Reader, stdout io.Writer, args []string) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}
	if _, err := fmt.Fprint(stdout, "GitHub repository (owner/repo or URL): "); err != nil {
		return "", err
	}
	reader := bufio.NewReader(stdin)
	input, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return "", errors.New("repository is required")
	}
	return input, nil
}

type scanOptions struct {
	json      bool
	compact   bool
	noColor   bool
	verbose   bool
	failUnder int
	timeout   time.Duration
}

func runScan(ctx context.Context, stdout io.Writer, input string, opts scanOptions) error {
	if opts.failUnder < 0 || opts.failUnder > 100 {
		return errors.New("--fail-under must be between 0 and 100")
	}
	ctx, cancel := context.WithTimeout(ctx, opts.timeout)
	defer cancel()

	ref, source, err := scanTarget(input, opts.timeout)
	if err != nil {
		return err
	}
	scan := scanner.New(source, config.Default(time.Now().UTC()), nil)
	start := time.Now()
	result, err := scan.Scan(ctx, ref)
	duration := time.Since(start)
	if err != nil {
		return err
	}
	if err := report.Write(stdout, result, report.Options{JSON: opts.json, Compact: opts.compact, NoColor: opts.noColor, Verbose: opts.verbose, Duration: duration}); err != nil {
		return err
	}
	if opts.failUnder > 0 && result.Score < opts.failUnder {
		return &scoreError{score: result.Score, threshold: opts.failUnder}
	}
	return nil
}

func scanTarget(input string, timeout time.Duration) (model.RepoRef, scanner.Source, error) {
	if local.IsPath(input) {
		ref := local.RefForPath(input)
		return ref, local.NewSource(input), nil
	}
	ref, err := model.ParseRepoRef(input)
	if err != nil {
		return model.RepoRef{}, nil, err
	}
	client := gh.NewClient(&http.Client{Timeout: timeout})
	return ref, client, nil
}

type scoreError struct {
	score     int
	threshold int
}

func (e *scoreError) Error() string {
	return fmt.Sprintf("repo health score %d is below threshold %d", e.score, e.threshold)
}

func exitCode(err error) int {
	var scoreErr *scoreError
	if errors.As(err, &scoreErr) {
		return 2
	}
	return 1
}
