package checks

import "strings"

func containsPath(files []string, candidates ...string) bool {
	for _, file := range files {
		normalized := strings.ToLower(strings.Trim(file, "/"))
		for _, candidate := range candidates {
			if normalized == strings.ToLower(strings.Trim(candidate, "/")) {
				return true
			}
		}
	}
	return false
}

func hasPrefix(files []string, prefixes ...string) bool {
	for _, file := range files {
		normalized := strings.ToLower(strings.Trim(file, "/"))
		for _, prefix := range prefixes {
			if strings.HasPrefix(normalized, strings.ToLower(strings.Trim(prefix, "/"))) {
				return true
			}
		}
	}
	return false
}

func countOlderThan[T any](items []T, older func(T) bool) int {
	total := 0
	for _, item := range items {
		if older(item) {
			total++
		}
	}
	return total
}
