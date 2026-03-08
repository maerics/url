package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  URLData
	}{
		{
			name:  "full url",
			input: "https://user:pass@example.com:8080/foo/bar?a=1&a=2&b=x#a%20b",
			want: URLData{
				Scheme: "https",
				Host: &HostData{
					String:   "user:pass@example.com:8080",
					User:     "user",
					Password: "pass",
					Name:     "example.com",
					Port:     "8080",
				},
				Path: &PathData{
					String: "/foo/bar",
					Parts:  []string{"foo", "bar"},
				},
				Query: &QueryData{
					String: "?a=1&a=2&b=x",
					Params: map[string]any{"a": []string{"1", "2"}, "b": "x"},
				},
				Fragment: &FragmentData{String: "a b", Raw: "#a%20b"},
			},
		},
		{
			name:  "scheme and host only",
			input: "https://example.com",
			want: URLData{
				Scheme: "https",
				Host:   &HostData{String: "example.com", Name: "example.com"},
			},
		},
		{
			name:  "path only",
			input: "/foo/bar/baz",
			want: URLData{
				Path: &PathData{String: "/foo/bar/baz", Parts: []string{"foo", "bar", "baz"}},
			},
		},
		{
			name:  "root path",
			input: "https://example.com/",
			want: URLData{
				Scheme: "https",
				Host:   &HostData{String: "example.com", Name: "example.com"},
				Path:   &PathData{String: "/"},
			},
		},
		{
			name:  "single query param",
			input: "https://example.com?q=hello",
			want: URLData{
				Scheme: "https",
				Host:   &HostData{String: "example.com", Name: "example.com"},
				Query:  &QueryData{String: "?q=hello", Params: map[string]any{"q": "hello"}},
			},
		},
		{
			name:  "plain fragment",
			input: "https://example.com#section",
			want: URLData{
				Scheme:   "https",
				Host:     &HostData{String: "example.com", Name: "example.com"},
				Fragment: &FragmentData{String: "section"},
			},
		},
		{
			name:  "non-canonical fragment encoding",
			input: "https://example.com#a%2Bb",
			want: URLData{
				Scheme:   "https",
				Host:     &HostData{String: "example.com", Name: "example.com"},
				Fragment: &FragmentData{String: "a+b", Raw: "#a%2Bb"},
			},
		},
		{
			name:  "user without password",
			input: "ftp://admin@files.example.com/pub",
			want: URLData{
				Scheme: "ftp",
				Host:   &HostData{String: "admin@files.example.com", User: "admin", Name: "files.example.com"},
				Path:   &PathData{String: "/pub", Parts: []string{"pub"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseURL(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.want, *got)
		})
	}
}
