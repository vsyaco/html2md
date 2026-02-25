# html2md

CLI utility for converting HTML pages (e.g. saved from Chrome) to Markdown.

## Install

```bash
# Homebrew
brew install vsyaco/tap/html2md

# Go
go install github.com/vsyaco/html2md@latest
```

Or download a binary from [Releases](https://github.com/vsyaco/html2md/releases).

## Usage

```
html2md [flags] <file.html>

Flags:
      --article          Extract main content only (readability)
      --domain string    Base URL for resolving relative links
  -h, --help             help for html2md
      --no-images        Remove images
  -o, --output string    Output .md file path
      --stdout           Print to stdout instead of file
  -v, --version          version for html2md
```

### Examples

```bash
# Convert to page.md next to the source file
html2md page.html

# Print to stdout
html2md page.html --stdout

# Extract main article content (strips nav, ads, etc.)
html2md page.html --article --stdout

# Write to a specific path
html2md page.html -o ~/notes/page.md

# Remove images and resolve relative links
html2md page.html --no-images --domain https://example.com
```

## Modes

**Default** — removes noisy elements (`nav`, `footer`, `header`, `aside`, forms, ads/sidebar/cookie/menu class elements) then converts the cleaned HTML to Markdown.

**Article** (`--article`) — uses [go-readability](https://github.com/go-shiori/go-readability) (port of Mozilla Readability.js) to extract the main article content before converting.

## License

MIT
