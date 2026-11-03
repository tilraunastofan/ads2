package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	dir, version, err := moduleDir() // in this repo, the main module itself
	if err != nil {
		t.Fatal(err)
	}
	css, err := generate(dir, version)
	if err != nil {
		t.Fatal(err)
	}
	d := filepath.ToSlash(dir)
	for _, want := range []string{
		`@import "` + d + `/assets/css/ads2.css";`,
		`@source "` + d + `/components";`,
		`@source "` + d + `/patterns";`,
	} {
		if !strings.Contains(css, want) {
			t.Errorf("missing %s in\n%s", want, css)
		}
	}
	// Every path it names exists, or Tailwind fails or silently generates nothing.
	for _, p := range []string{"assets/css/ads2.css", "components", "patterns", "assets/js"} {
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			t.Errorf("%s: %v", p, err)
		}
	}
}

func TestGenerateRefusesAnOldModule(t *testing.T) {
	if _, err := generate(t.TempDir(), "v0.0.1"); err == nil || !strings.Contains(err.Error(), "go get") {
		t.Errorf("want an error that says to update, got %v", err)
	}
}
