package model

import "testing"

func TestParseRepoRef(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    RepoRef
		wantErr bool
	}{
		{name: "short", input: "openai/codex", want: RepoRef{Owner: "openai", Name: "codex"}},
		{name: "url", input: "https://github.com/openai/codex", want: RepoRef{Owner: "openai", Name: "codex"}},
		{name: "git suffix", input: "https://github.com/openai/codex.git", want: RepoRef{Owner: "openai", Name: "codex"}},
		{name: "wrong host", input: "https://example.com/openai/codex", wantErr: true},
		{name: "http", input: "http://github.com/openai/codex", wantErr: true},
		{name: "too many parts", input: "openai/codex/issues", wantErr: true},
		{name: "bad chars", input: "openai/co dex", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRepoRef(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
