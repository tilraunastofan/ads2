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
	// skipAt skips one match, given the line and where the match starts (for a
	// context the regexp can't express: RE2 has no lookbehind).
	skipAt func(line string, start int) bool
}

const palette = `slate|gray|zinc|neutral|stone|red|orange|amber|yellow|lime|green|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose`

// Colour utilities, with the side and axis forms (border-t-, divide-y-, ring-offset-).
const colourUtils = `bg|text|border(?:-[xytrblse])?|ring(?:-offset)?|outline|fill|stroke|from|via|to|divide(?:-[xy])?|decoration|shadow|caret|accent|placeholder`

// fragmentRef is an attribute value that names an element, not a colour: href="#add",
// hx-target="#bad", target:#fab. A hex-looking id there is not a hex colour.
var fragmentRef = regexp.MustCompile(`(?:href|hx-[a-z:-]+|for|aria-[a-z]+|data-[a-z-]+)="$|target:$`)

var rules = []rule{
	{
		name: "hex colour",
		re:   regexp.MustCompile(`#(?:[0-9a-fA-F]{8}|[0-9a-fA-F]{6}|[0-9a-fA-F]{3,4})\b`),
		why:  "colours come from tokens (bg-primary, text-muted-foreground); change tokens/tokens.json instead",
		skipAt: func(line string, start int) bool {
			return fragmentRef.MatchString(line[:start])
		},
	},
	{
		name: "Tailwind palette colour",
		re:   regexp.MustCompile(`\b(?:` + colourUtils + `)-(?:(?:` + palette + `)-\d{2,3}|white|black)\b`),
		why:  "use a token utility (bg-primary, text-positive, bg-ink); on a fill use its *-foreground, never white",
	},
	{
		name: "Tailwind text size",
		re:   regexp.MustCompile(`(?:^|[\s"':])text-(?:xs|sm|base|lg|xl|[2-9]xl)\b`),
		why:  "use the type scale (text-caption, text-body-sm, text-body, text-display-sm …); the sizes match, the tokens carry line height and weight",
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
		name: "bool aria attribute",
		re:   regexp.MustCompile(`"aria-[a-z]+"\s*:\s*[^",}]*(?:!=|==|\btrue\b|\bfalse\b)`),
		why:  `templ renders a Go bool as a bare attribute (aria-invalid), which ARIA reads as false; pass "true" as a string, or leave the attribute out`,
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
			if r.allow != nil && r.allow(line) {
				continue
			}
			for _, loc := range r.re.FindAllStringIndex(line, -1) {
				if r.skipAt == nil || !r.skipAt(line, loc[0]) {
					out = append(out, violation{line: i + 1, rule: r, text: strings.TrimSpace(line)})
					break
				}
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
		`<div class="bg-[#005ab4]">`:                                 "hex colour",
		`<span class="text-blue-600">`:                               "Tailwind palette colour",
		`@button.Button(button.Props{Class: "text-white"})`:          "Tailwind palette colour",
		`<p class="text-[40px]">`:                                    "arbitrary font size",
		`<a hx-delete="/x" hx-confirm="Sure?">`:                      "hx-confirm",
		`<button type="submit">`:                                     "raw form or table element",
		`<input type="text" name="q">`:                               "raw form or table element",
		`<h1>Releases</h1>`:                                          "page title outside PageHeader",
		`Status: "ocp-prod-1 · 42 applications"`:                     "middle-dot meta string",
		`Attributes: templ.Attributes{"aria-invalid": e.Name != ""}`: "bool aria attribute",
		`templ.Attributes{"aria-expanded": true}`:                    "bool aria attribute",
		`<div class="border-t-red-500">`:                             "Tailwind palette colour",
		`<div class="border-x-blue-600">`:                            "Tailwind palette colour",
		`<div class="divide-y-slate-200">`:                           "Tailwind palette colour",
		`<div class="ring-offset-white">`:                            "Tailwind palette colour",
		`<p class="mt-1 text-sm text-muted-foreground">`:             "Tailwind text size",
		`@card.Title(card.TitleProps{Class: "text-xs"})`:             "Tailwind text size",
		`<a href="#add" class="fill-[#fab]">`:                        "hex colour",
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
		`templ.Attributes{"aria-invalid": "true", "aria-hidden": "true"}`,
		`<a href="#add">`,
		`<div hx-target="#bad" hx-status:5xx="target:#fab">`,
		`<p class="text-body-sm text-caption text-muted-foreground">`,
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

// TestDocExamplesFollowDesignRules runs the same rules over the ```templ examples
// in docs/: teams copy those, so a middot or a raw <button> there spreads as
// surely as one in a template.
func TestDocExamplesFollowDesignRules(t *testing.T) {
	block := regexp.MustCompile("(?s)```templ\n(.*?)```")
	n := 0
	err := filepath.WalkDir(filepath.Join("..", "docs"), func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".md") {
			return err
		}
		src, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel("..", p)
		for _, m := range block.FindAllStringSubmatch(string(src), -1) {
			n++
			for _, v := range check("docs/", m[1]) {
				t.Errorf("%s: example breaks %s: %s\n\t%s", filepath.ToSlash(rel), v.rule.name, v.rule.why, v.text)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatal("no ```templ examples found in docs/")
	}
}
