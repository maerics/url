# url

Parse URLs and print structured data in JSON or YAML.

## Install

```sh
# Homebrew
brew install maerics/datautils/url

# Go
go install github.com/maerics/url@latest
```

## Usage

```
Usage: url <urls> ... [flags]

Arguments:
  <urls> ...    URLs to parse

Flags:
  -h, --help           Show context-sensitive help.
  -m, --mode="yaml"    Output mode (json or yaml)
```

## Examples

Parse a full URL (default YAML output):

```sh
$ url 'https://user:pass@example.com:8080/a/b?x=1&y=2#frag'
scheme: https
host:
  string: user:pass@example.com:8080
  user: user
  password: pass
  name: example.com
  port: "8080"
path:
  string: /a/b
  parts:
    - a
    - b
query:
  string: ?x=1&y=2
  params:
    x: "1"
    y: "2"
fragment:
  string: frag
```

JSON output with `-m json`:

```sh
$ url -m json 'https://example.com/search?q=hello'
{
  "scheme": "https",
  "host": {
    "string": "example.com",
    "name": "example.com"
  },
  "path": {
    "string": "/search",
    "parts": [
      "search"
    ]
  },
  "query": {
    "string": "?q=hello",
    "params": {
      "q": "hello"
    }
  }
}
```

Parse multiple URLs (YAML documents separated by `---`):

```sh
$ url 'https://example.com' 'http://localhost:3000/api'
scheme: https
host:
  string: example.com
  name: example.com
---
scheme: http
host:
  string: localhost:3000
  name: localhost
  port: "3000"
path:
  string: /api
  parts:
    - api
```

If `jq` (JSON mode) or `yq` (YAML mode) is installed, output is piped through it for syntax highlighting.

## License

MIT
