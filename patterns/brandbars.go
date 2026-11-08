package patterns

import "github.com/a-h/templ"

// barAttrs marks the top bar ("top") or the sidebar ("side") for the brand
// gradient rules in assets/css/ads2.css, or returns nil when brand bars are off,
// so a shell without them renders exactly as before.
func barAttrs(on bool, which string) templ.Attributes {
	if !on {
		return nil
	}
	return templ.Attributes{"data-brand-bar": which}
}
