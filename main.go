package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"
)

var cli struct {
	Mode string   `help:"Output mode (json or yaml)" short:"m" default:"yaml" enum:"json,yaml"`
	Urls []string `arg:"" help:"URLs to parse"`
}

type HostData struct {
	String   string `json:"string,omitempty" yaml:"string,omitempty"`
	User     string `json:"user,omitempty" yaml:"user,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	Hostname string `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Port     string `json:"port,omitempty" yaml:"port,omitempty"`
}

type PathData struct {
	String string   `json:"string,omitempty" yaml:"string,omitempty"`
	Parts  []string `json:"parts,omitempty" yaml:"parts,omitempty"`
}

type QueryData struct {
	String string         `json:"string,omitempty" yaml:"string,omitempty"`
	Params map[string]any `json:"params,omitempty" yaml:"params,omitempty"`
}

type FragmentData struct {
	String string `json:"string,omitempty" yaml:"string,omitempty"`
	Raw    string `json:"raw,omitempty" yaml:"raw,omitempty"`
}

type URLData struct {
	Scheme   string        `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	Host     *HostData     `json:"host,omitempty" yaml:"host,omitempty"`
	Path     *PathData     `json:"path,omitempty" yaml:"path,omitempty"`
	Query    *QueryData    `json:"query,omitempty" yaml:"query,omitempty"`
	Fragment *FragmentData `json:"fragment,omitempty" yaml:"fragment,omitempty"`
}

func parseURL(rawURL string) (*URLData, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	data := &URLData{}

	if u.Scheme != "" {
		data.Scheme = u.Scheme
	}

	if u.Host != "" || u.User != nil {
		hd := &HostData{}
		authority := u.Host
		if u.User != nil {
			hd.User = u.User.Username()
			if p, ok := u.User.Password(); ok {
				hd.Password = p
			}
			authority = u.User.String() + "@" + u.Host
		}
		hd.String = authority
		hd.Hostname = u.Hostname()
		hd.Port = u.Port()
		data.Host = hd
	}

	if u.Path != "" {
		var parts []string
		for _, p := range strings.Split(strings.Trim(u.Path, "/"), "/") {
			if p != "" {
				parts = append(parts, p)
			}
		}
		data.Path = &PathData{
			String: u.Path,
			Parts:  parts,
		}
	}

	if u.RawQuery != "" {
		values, _ := url.ParseQuery(u.RawQuery)
		params := make(map[string]any, len(values))
		for k, v := range values {
			if len(v) == 1 {
				params[k] = v[0]
			} else {
				params[k] = v
			}
		}
		data.Query = &QueryData{
			String: "?" + u.RawQuery,
			Params: params,
		}
	}

	if u.Fragment != "" {
		fd := &FragmentData{String: u.Fragment}
		if esc := u.EscapedFragment(); esc != u.Fragment {
			fd.Raw = "#" + esc
		}
		data.Fragment = fd
	}

	return data, nil
}

// prettyWriter returns a writer connected to jq/yq for pretty-printing,
// falling back to os.Stdout if the command is unavailable or fails to start.
func prettyWriter(mode string) (io.Writer, func()) {
	var name string
	switch mode {
	case "json":
		name = "jq"
	case "yaml":
		name = "yq"
	}

	if path, err := exec.LookPath(name); err == nil {
		cmd := exec.Command(path, ".")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if pipe, err := cmd.StdinPipe(); err == nil {
			if err = cmd.Start(); err == nil {
				return pipe, func() { pipe.Close(); cmd.Wait() }
			}
		}
	}

	return os.Stdout, func() {}
}

func main() {
	kong.Parse(&cli)

	w, cleanup := prettyWriter(cli.Mode)
	defer cleanup()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	for i, rawURL := range cli.Urls {
		data, err := parseURL(rawURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		switch cli.Mode {
		case "json":
			if err = enc.Encode(data); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		case "yaml":
			if i > 0 {
				fmt.Fprintln(w, "---")
			}
			out, err := yaml.Marshal(data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			w.Write(out)
		}
	}
}
