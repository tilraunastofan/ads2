package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// pair is a contrast requirement between two color tokens. When tint > 0 the
// background is fg drawn at that opacity over bg (a "soft" destructive button).
type pair struct {
	fg, bg string
	min    float64
	tint   float64
	only   string // restrict to one theme id; empty checks every theme
}

const (
	text = 4.5 // WCAG 1.4.3 normal text
	ui   = 3.0 // WCAG 1.4.11 control boundaries, focus rings, chart marks
)

// requirements are the contrast claims docs/design-system.md makes. A failing
// pair fails `go run ./cmd/gentheme` and `go test ./cmd/gentheme`.
var requirements = []pair{
	{"foreground", "background", text, 0, ""}, {"foreground", "card", text, 0, ""}, {"foreground", "muted", text, 0, ""},
	{"card-foreground", "card", text, 0, ""}, {"popover-foreground", "popover", text, 0, ""},
	{"primary-foreground", "primary", text, 0, ""}, {"secondary-foreground", "secondary", text, 0, ""},
	{"destructive-foreground", "destructive", text, 0, ""}, {"positive-foreground", "positive", text, 0, ""},
	{"warning-foreground", "warning", text, 0, ""}, {"accent-foreground", "accent", text, 0, ""},
	{"sidebar-foreground", "sidebar", text, 0, ""}, {"sidebar-accent-foreground", "sidebar-accent", text, 0, ""},
	{"sidebar-primary-foreground", "sidebar-primary", text, 0, ""}, {"ink-foreground", "ink", text, 0, ""},
	{"ink-foreground", "ink-hover", text, 0, ""}, {"selection-foreground", "selection", text, 0, ""},
	{"surface-foreground", "surface", text, 0, ""},

	{"muted-foreground", "background", text, 0, ""}, {"muted-foreground", "card", text, 0, ""},
	{"muted-foreground", "muted", text, 0, ""}, {"muted-foreground", "secondary", text, 0, ""},
	{"muted-foreground", "accent", text, 0, ""}, {"muted-foreground", "surface", text, 0, ""},
	{"primary", "background", text, 0, ""}, {"primary", "card", text, 0, ""}, {"primary", "muted", text, 0, ""},
	{"destructive", "background", text, 0, ""}, {"destructive", "card", text, 0, ""}, {"destructive", "muted", text, 0, ""},
	{"destructive", "background", text, 0.1, ""}, {"destructive", "card", text, 0.1, ""},
	{"destructive", "card", text, 0.2, "dark"}, {"positive", "card", text, 0.2, "dark"}, {"warning", "card", text, 0.2, "dark"},
	{"positive", "background", text, 0, ""}, {"positive", "card", text, 0, ""}, {"positive", "muted", text, 0, ""},
	{"positive", "background", text, 0.1, ""}, {"positive", "card", text, 0.1, ""},
	{"warning", "background", text, 0, ""}, {"warning", "card", text, 0, ""}, {"warning", "muted", text, 0, ""},
	{"warning", "background", text, 0.1, ""}, {"warning", "card", text, 0.1, ""},
	{"gold", "background", text, 0, ""}, {"gold", "card", text, 0, ""},
	{"plum", "card", text, 0, ""}, {"slate", "card", text, 0, ""},

	{"input", "background", ui, 0, ""}, {"input", "card", ui, 0, ""},
	{"ring", "background", ui, 0, ""}, {"ring", "card", ui, 0, ""},
	{"chart-1", "card", ui, 0, ""}, {"chart-2", "card", ui, 0, ""}, {"chart-3", "card", ui, 0, ""},
	{"chart-4", "card", ui, 0, ""}, {"chart-5", "card", ui, 0, ""},
}

func parseHex(s string) ([3]float64, error) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(s) != 6 {
		return [3]float64{}, fmt.Errorf("%q is not a #rrggbb color", s)
	}
	var c [3]float64
	for i := range c {
		v, err := strconv.ParseUint(s[i*2:i*2+2], 16, 8)
		if err != nil {
			return c, err
		}
		c[i] = float64(v)
	}
	return c, nil
}

func luminance(c [3]float64) float64 {
	var l [3]float64
	for i, v := range c {
		v /= 255
		if v <= 0.04045 {
			l[i] = v / 12.92
		} else {
			l[i] = math.Pow((v+0.055)/1.055, 2.4)
		}
	}
	return 0.2126*l[0] + 0.7152*l[1] + 0.0722*l[2]
}

func ratio(a, b [3]float64) float64 {
	x, y := luminance(a), luminance(b)
	if x < y {
		x, y = y, x
	}
	return (x + 0.05) / (y + 0.05)
}

// checkContrast returns one message per requirement that fails in any theme.
func checkContrast(tk *tokens) []string {
	var fails []string
	for _, th := range tk.Color.Themes {
		get := func(name string) ([3]float64, error) {
			for _, t := range tk.Color.Tokens {
				if t.Name == name {
					v, err := value(t, th.ID)
					if err != nil {
						return [3]float64{}, err
					}
					return parseHex(v)
				}
			}
			return [3]float64{}, fmt.Errorf("unknown token %q", name)
		}
		for _, p := range requirements {
			if p.only != "" && p.only != th.ID {
				continue
			}
			fg, err1 := get(p.fg)
			bg, err2 := get(p.bg)
			if err1 != nil || err2 != nil {
				fails = append(fails, fmt.Sprintf("%s: %s on %s: %v %v", th.ID, p.fg, p.bg, err1, err2))
				continue
			}
			label := p.bg
			if p.tint > 0 {
				for i := range bg {
					bg[i] = fg[i]*p.tint + bg[i]*(1-p.tint)
				}
				label = fmt.Sprintf("%.0f%% %s over %s", p.tint*100, p.fg, p.bg)
			}
			if r := ratio(fg, bg); r < p.min {
				fails = append(fails, fmt.Sprintf("%s: %s on %s is %.2f:1, needs %.1f:1", th.ID, p.fg, label, r, p.min))
			}
		}
	}
	return fails
}
