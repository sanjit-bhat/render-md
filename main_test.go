package main

import (
	"os"
	"strings"
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	input, err := os.ReadFile("testdata/input.md")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/want.html")
	if err != nil {
		t.Fatal(err)
	}

	got, err := renderMarkdown(string(input))
	if err != nil {
		t.Fatalf("renderMarkdown: %v", err)
	}
	// API response includes a trailing newline.
	got = strings.TrimRight(got, "\n")
	if got != strings.TrimRight(string(want), "\n") {
		t.Errorf("output mismatch\ngot:\n%s\n\nwant:\n%s", got, want)
	}
}
