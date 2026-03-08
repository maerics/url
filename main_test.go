package main

import (
	"testing"
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
					Hostname: "example.com",
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
				Host:   &HostData{String: "example.com", Hostname: "example.com"},
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
				Host:   &HostData{String: "example.com", Hostname: "example.com"},
				Path:   &PathData{String: "/"},
			},
		},
		{
			name:  "single query param",
			input: "https://example.com?q=hello",
			want: URLData{
				Scheme: "https",
				Host:   &HostData{String: "example.com", Hostname: "example.com"},
				Query:  &QueryData{String: "?q=hello", Params: map[string]any{"q": "hello"}},
			},
		},
		{
			name:  "plain fragment",
			input: "https://example.com#section",
			want: URLData{
				Scheme:   "https",
				Host:     &HostData{String: "example.com", Hostname: "example.com"},
				Fragment: &FragmentData{String: "section"},
			},
		},
		{
			name:  "non-canonical fragment encoding",
			input: "https://example.com#a%2Bb",
			want: URLData{
				Scheme:   "https",
				Host:     &HostData{String: "example.com", Hostname: "example.com"},
				Fragment: &FragmentData{String: "a+b", Raw: "#a%2Bb"},
			},
		},
		{
			name:  "user without password",
			input: "ftp://admin@files.example.com/pub",
			want: URLData{
				Scheme: "ftp",
				Host:   &HostData{String: "admin@files.example.com", User: "admin", Hostname: "files.example.com"},
				Path:   &PathData{String: "/pub", Parts: []string{"pub"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseURL(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("got nil result")
			}
			assertEqual(t, "Scheme", tt.want.Scheme, got.Scheme)
			assertHost(t, tt.want.Host, got.Host)
			assertPath(t, tt.want.Path, got.Path)
			assertQuery(t, tt.want.Query, got.Query)
			assertFragment(t, tt.want.Fragment, got.Fragment)
		})
	}
}

func assertEqual(t *testing.T, field, want, got string) {
	t.Helper()
	if want != got {
		t.Errorf("%s: want %q, got %q", field, want, got)
	}
}

func assertHost(t *testing.T, want, got *HostData) {
	t.Helper()
	if want == nil && got == nil {
		return
	}
	if want == nil || got == nil {
		t.Errorf("host: want %v, got %v", want, got)
		return
	}
	assertEqual(t, "host.string", want.String, got.String)
	assertEqual(t, "host.user", want.User, got.User)
	assertEqual(t, "host.password", want.Password, got.Password)
	assertEqual(t, "host.hostname", want.Hostname, got.Hostname)
	assertEqual(t, "host.port", want.Port, got.Port)
}

func assertPath(t *testing.T, want, got *PathData) {
	t.Helper()
	if want == nil && got == nil {
		return
	}
	if want == nil || got == nil {
		t.Errorf("path: want %v, got %v", want, got)
		return
	}
	assertEqual(t, "path.string", want.String, got.String)
	if len(want.Parts) != len(got.Parts) {
		t.Errorf("path.parts: want %v, got %v", want.Parts, got.Parts)
		return
	}
	for i := range want.Parts {
		assertEqual(t, "path.parts["+string(rune('0'+i))+"]", want.Parts[i], got.Parts[i])
	}
}

func assertQuery(t *testing.T, want, got *QueryData) {
	t.Helper()
	if want == nil && got == nil {
		return
	}
	if want == nil || got == nil {
		t.Errorf("query: want %v, got %v", want, got)
		return
	}
	assertEqual(t, "query.string", want.String, got.String)
	for k, wv := range want.Params {
		gv, ok := got.Params[k]
		if !ok {
			t.Errorf("query.params[%q]: missing", k)
			continue
		}
		switch wval := wv.(type) {
		case string:
			if gval, ok := gv.(string); !ok || gval != wval {
				t.Errorf("query.params[%q]: want %q, got %v", k, wval, gv)
			}
		case []string:
			gval, ok := gv.([]string)
			if !ok {
				t.Errorf("query.params[%q]: want []string, got %T", k, gv)
				continue
			}
			if len(gval) != len(wval) {
				t.Errorf("query.params[%q]: want %v, got %v", k, wval, gval)
				continue
			}
			for i := range wval {
				if wval[i] != gval[i] {
					t.Errorf("query.params[%q][%d]: want %q, got %q", k, i, wval[i], gval[i])
				}
			}
		}
	}
}

func assertFragment(t *testing.T, want, got *FragmentData) {
	t.Helper()
	if want == nil && got == nil {
		return
	}
	if want == nil || got == nil {
		t.Errorf("fragment: want %v, got %v", want, got)
		return
	}
	assertEqual(t, "fragment.string", want.String, got.String)
	assertEqual(t, "fragment.raw", want.Raw, got.Raw)
}
