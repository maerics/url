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
  -f, --format="json"    Output format (json or yaml)
```

## Examples

Parse a full URL (default JSON output):

```sh
$ url 'https://user:pass@example.com:8080/a/b?x=1&y=2#frag'
{
  "scheme": "https",
  "host": {
    "authority": "user:pass@example.com:8080",
    "user": "user",
    "password": "pass",
    "name": "example.com",
    "port": "8080"
  },
  "path": {
    "string": "/a/b",
    "parts": [
      "a",
      "b"
    ]
  },
  "query": {
    "string": "?x=1&y=2",
    "params": {
      "x": "1",
      "y": "2"
    }
  },
  "fragment": {
    "string": "frag"
  }
}
```

YAML output with `-f yaml`:

```sh
$ url -f yaml 'https://example.com/search?q=hello'
scheme: https
host:
  name: example.com
path:
  string: /search
  parts:
    - search
query:
  string: ?q=hello
  params:
    q: hello
```

Parse multiple URLs (JSON documents printed one per line):

```sh
$ url 'https://example.com' 'http://localhost:3000/api'
{"scheme":"https","host":{"name":"example.com"}}
{"scheme":"http","host":{"name":"localhost","port":"3000"},"path":{"string":"/api","parts":["api"]}}
```

If `jq` (JSON mode) or `yq` (YAML mode) is installed, output is piped through it for syntax highlighting.

## License

MIT
