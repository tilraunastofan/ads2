package utils

import (
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func TestCNKeepsTypeScaleAndColour(t *testing.T) {
	tests := []struct {
		in   []any
		want string
	}{
		{[]any{"font-mono text-code font-normal text-muted-foreground"}, "font-mono text-code font-normal text-muted-foreground"},
		{[]any{"p-2 whitespace-nowrap", "text-caption text-muted-foreground"}, "p-2 whitespace-nowrap text-caption text-muted-foreground"},
		{[]any{"text-muted-foreground", "text-caption"}, "text-muted-foreground text-caption"},
		// Sizes still replace each other: the last one wins.
		{[]any{"text-sm", "text-body-sm"}, "text-body-sm"},
		{[]any{"text-body-sm", "text-heading-3"}, "text-heading-3"},
		// Colours still replace each other.
		{[]any{"text-foreground", "text-destructive"}, "text-destructive"},
	}
	for _, tt := range tests {
		if got := CN(tt.in...); got != tt.want {
			t.Errorf("CN(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// The list must name every size the theme defines.
func TestTypeScaleMatchesTheme(t *testing.T) {
	css, err := os.ReadFile("../assets/css/arion-theme.css")
	if err != nil {
		t.Fatal(err)
	}
	var theme []string
	for _, m := range regexp.MustCompile(`(?m)^\s+--text-([a-z0-9-]+?):`).FindAllStringSubmatch(string(css), -1) {
		if !strings.Contains(m[1], "--") && !slices.Contains(theme, m[1]) { // not text-body--line-height
			theme = append(theme, m[1])
		}
	}
	slices.Sort(theme)
	got := slices.Sorted(slices.Values(typeScale))
	if !slices.Equal(got, theme) {
		t.Errorf("typeScale = %v\ntheme sizes = %v", got, theme)
	}
}
