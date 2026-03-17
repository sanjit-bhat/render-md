package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const htmlTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>%s</title>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.8.1/github-markdown.min.css">
  <style>
    body { box-sizing: border-box; min-width: 200px; max-width: 980px; margin: 0 auto; padding: 45px; }
  </style>
</head>
<body class="markdown-body">
%s
</body>
</html>
`

func main() {
	input, output, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	contents, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", input, err)
		os.Exit(1)
	}

	fragment, err := renderMarkdown(string(contents))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error rendering markdown: %v\n", err)
		os.Exit(1)
	}

	title := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
	html := fmt.Sprintf(htmlTemplate, title, fragment)

	if err := os.WriteFile(output, []byte(html), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing %s: %v\n", output, err)
		os.Exit(1)
	}

	fmt.Println(output)
}

func parseArgs(args []string) (input, output string, err error) {
	fs := flag.NewFlagSet("render-md", flag.ContinueOnError)
	fs.StringVar(&output, "output", "", "output file path (default: input with .html extension)")
	fs.StringVar(&output, "o", "", "output file path (shorthand)")
	if err := fs.Parse(args); err != nil {
		return "", "", err
	}
	if fs.NArg() != 1 {
		return "", "", fmt.Errorf("usage: render-md <file.md> [-output <file.html>]")
	}
	input = fs.Arg(0)
	if filepath.Ext(input) != ".md" {
		return "", "", fmt.Errorf("input file must have a .md extension")
	}
	if _, err := os.Stat(input); os.IsNotExist(err) {
		return "", "", fmt.Errorf("file not found: %s", input)
	}
	if output == "" {
		output = strings.TrimSuffix(input, ".md") + ".html"
	}
	return input, output, nil
}

func renderMarkdown(text string) (string, error) {
	cmd := exec.Command("gh", "api", "/markdown",
		"--method", "POST",
		"--field", "text="+text,
		"--field", "mode=markdown",
		"--header", "Accept: application/vnd.github+json",
		"--header", "X-GitHub-Api-Version: 2026-03-10",
	)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("gh api failed: %s", exitErr.Stderr)
		}
		return "", err
	}
	return string(out), nil
}
