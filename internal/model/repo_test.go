package model

import "testing"

func TestParseRepoRef(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    RepoRef
		wantErr bool
	}{
		{name: "short", input: "github/cli", want: RepoRef{Owner: "github", Name: "cli"}},
		{name: "url", input: "https://github.com/github/cli", want: RepoRef{Owner: "github", Name: "cli"}},
		{name: "git suffix", input: "https://github.com/github/cli.git", want: RepoRef{Owner: "github", Name: "cli"}},
		{name: "wrong host", input: "https://example.com/github/cli", wantErr: true},
		{name: "http", input: "http://github.com/github/cli", wantErr: true},
		{name: "too many parts", input: "github/cli/issues", wantErr: true},
		{name: "bad chars", input: "github/cl i", wantErr: true},
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
