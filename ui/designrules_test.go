// Package ui holds the showcase's layouts and pages. This test enforces the "Rules
// for UI work" in CLAUDE.md on every .templ file in ui/ and patterns/, so a hex
// colour or a raw <button> fails `go test` instead of waiting for review.
package ui

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

type rule struct {
	name string
	re   *regexp.Regexp
	why  string
	// allow skips a match on this line (e.g. <input type="hidden">).
	allow func(line string) bool
	// only applies the rule to files under this directory (slash-separated); empty means all.
	only string
}

const palette = `slate|gray|zinc|neutral|stone|red|orange|amber|yellow|lime|green|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose`
const colourUtils = `bg|text|border|ring|outline|fill|stroke|from|via|to|divide|decoration|shadow|caret|accent|placeholder`

var rules = []rule{
	{
		name: "hex colour",
		re:   regexp.MustCompile(`#(?:[0-9a-fA-F]{8}|[0-9a-fA-F]{6}|[0-9a-fA-F]{3})\b`),
		why:  "colours come from tokens (bg-primary, text-muted-foreground); change tokens/tokens.json instead",
	},
	{
		name: "Tailwind palette colour",
		re:   regexp.MustCompile(`\b(?:` + colourUtils + `)-(?:(?:` + palette + `)-\d{2,3}|white|black)\b`),
		why:  "use a token utility (bg-primary, text-positive, bg-ink); on a fill use its *-foreground, never white",
	},
	{
		name: "arbitrary font size",
		re:   regexp.MustCompile(`\btext-\[(?:\d|length:)`),
		why:  "use the type scale (text-display-sm, text-body-sm, text-code)",
	},
	{
		name: "hx-confirm",
		re:   regexp.MustCompile(`\bhx-confirm\b`),
		why:  "destructive actions go through alertdialog",
	},
	{
		name: "raw form or table element",
		re:   regexp.MustCompile(`<(?:button|table|select|textarea|dialog|input)\b`),
		why:  "use the shadcn-templ component (button.Button, table.Table, nativeselect …)",
		allow: func(line string) bool {
			return strings.HasPrefix(strings.TrimSpace(line), "<input") && strings.Contains(line, `type="hidden"`)
		},
	},
	{
		name: "middle-dot meta string",
		re:   regexp.MustCompile(`·`),
		why:  "don't join metadata with middle dots; write a sentence or use separate elements",
	},
	{
		name: "page title outside PageHeader",
		re:   regexp.MustCompile(`<h1\b`),
		why:  "the one display title per page comes from patterns.PageHeader",
		only: "ui/",
	},
}

type violation struct {
	line int
	rule rule
	text string
}

// check returns the rule violations in one .templ source. Go comments are skipped,
// so a rule can be quoted in a doc comment.
func check(path, src string) []violation {
	var out []violation
	for i, line := range strings.Split(src, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}
		for _, r := range rules {
			if r.only != "" && !strings.HasPrefix(path, r.only) {
				continue
			}
			if r.re.MatchString(line) && (r.allow == nil || !r.allow(line)) {
				out = append(out, violation{line: i + 1, rule: r, text: strings.TrimSpace(line)})
			}
		}
	}
	return out
}

func TestTemplatesFollowDesignRules(t *testing.T) {
	root := ".."
	n := 0
	for _, dir := range []string{"ui", "patterns"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(p, ".templ") {
				return err
			}
			src, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(root, p)
			rel = filepath.ToSlash(rel)
			n++
			for _, v := range check(rel, string(src)) {
				t.Errorf("%s:%d: %s: %s\n\t%s", rel, v.line, v.rule.name, v.rule.why, v.text)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if n == 0 {
		t.Fatal("no .templ files found; is the test running from ui/?")
	}
}

// TestRulesCatchViolations keeps the checker honest: each sample must trip exactly
// the named rule, and the allowed forms must pass.
func TestRulesCatchViolations(t *testing.T) {
	bad := map[string]string{
		`<div class="bg-[#005ab4]">`:                        "hex colour",
		`<span class="text-blue-600">`:                      "Tailwind palette colour",
		`@button.Button(button.Props{Class: "text-white"})`: "Tailwind palette colour",
		`<p class="text-[40px]">`:                           "arbitrary font size",
		`<a hx-delete="/x" hx-confirm="Sure?">`:             "hx-confirm",
		`<button type="submit">`:                            "raw form or table element",
		`<input type="text" name="q">`:                      "raw form or table element",
		`<h1>Releases</h1>`:                                 "page title outside PageHeader",
		`Status: "ocp-prod-1 · 42 applications"`:            "middle-dot meta string",
	}
	for src, want := range bad {
		vs := check("ui/pages/x.templ", src)
		if len(vs) != 1 || vs[0].rule.name != want {
			t.Errorf("%s: got %v, want one %q", src, vs, want)
		}
	}
	good := []string{
		`<div class="bg-primary text-primary-foreground hover:bg-primary/90">`,
		`<input type="hidden" name="from" value="2026-10-01">`,
		`hx-target="#app-rows"`,
		`<div class="w-[260px] h-[240px]">`,
		`// never text-blue-600 or <button>`,
	}
	for _, src := range good {
		if vs := check("ui/pages/x.templ", src); len(vs) != 0 {
			t.Errorf("%s: unexpected %v", src, vs)
		}
	}
	if vs := check("patterns/pageheader.templ", `<h1 class="font-display text-display-sm">`); len(vs) != 0 {
		t.Errorf("PageHeader may render the h1: %v", vs)
	}
}
