package utils

import (
	"slices"

	twmerge "github.com/Oudwins/tailwind-merge-go/pkg/twmerge"
)

// typeScale is the type scale's sizes (text-<name>), from tokens/tokens.json via
// assets/css/arion-theme.css. TestTypeScaleMatchesTheme fails when they drift.
var typeScale = []string{
	"display-xl", "display-lg", "display-md", "display-sm",
	"heading-1", "heading-2", "heading-3", "heading-4",
	"body-lg", "body", "body-sm", "label", "caption", "figure", "code",
}

// merge is tailwind-merge that knows the type scale is font sizes. Out of the box
// it reads text-caption or text-code as a text colour (anything after text- that
// isn't a t-shirt size is one), so CN("text-code text-muted-foreground") dropped
// the size and CN("text-muted-foreground text-caption") the colour.
var merge = func() twmerge.TwMergeFn {
	cfg := twmerge.MakeDefaultConfig()
	text := cfg.ClassGroups.NextPart["text"]
	text.Validators = append([]twmerge.ClassGroupValidator{{
		Fn:           func(v string) bool { return slices.Contains(typeScale, v) },
		ClassGroupId: "font-size",
	}}, text.Validators...)
	cfg.ClassGroups.NextPart["text"] = text
	return twmerge.CreateTwMerge(cfg, nil)
}()
