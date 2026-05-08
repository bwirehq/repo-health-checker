package github

import "strings"

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
