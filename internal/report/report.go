package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/bwirehq/repo-health-checker/internal/model"
)

type Options struct {
	JSON    bool
	NoColor bool
	Verbose bool
}

func Write(w io.Writer, result model.ScanResult, opts Options) error {
	if opts.JSON {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}
	return writeText(w, result, opts)
}

func writeText(w io.Writer, result model.ScanResult, opts Options) error {
	_, err := fmt.Fprintf(w, "Repo Health: %d/%d (%s)\n\n", result.Score, result.MaxScore, result.Grade)
	if err != nil {
		return err
	}
	for _, check := range result.Checks {
		line := fmt.Sprintf("%s %s", marker(check.Status, opts.NoColor), check.Title)
		if opts.Verbose {
			line = fmt.Sprintf("%s %d/%d", line, check.Points, check.MaxPoints)
		}
		if _, err := fmt.Fprintf(w, "%s\n  %s\n", line, check.Summary); err != nil {
			return err
		}
	}

	if len(result.Recommendations) > 0 {
		if _, err := fmt.Fprintln(w, "\nTop fixes:"); err != nil {
			return err
		}
		for i, rec := range result.Recommendations {
			if _, err := fmt.Fprintf(w, "%d. %s - %s\n", i+1, rec.Title, rec.Detail); err != nil {
				return err
			}
		}
	}
	return nil
}

func marker(status model.Status, noColor bool) string {
	text := strings.ToUpper(string(status))
	if noColor {
		return text
	}
	switch status {
	case model.StatusPass:
		return "\x1b[32mPASS\x1b[0m"
	case model.StatusWarn:
		return "\x1b[33mWARN\x1b[0m"
	case model.StatusFail:
		return "\x1b[31mFAIL\x1b[0m"
	default:
		return "\x1b[36mINFO\x1b[0m"
	}
}
