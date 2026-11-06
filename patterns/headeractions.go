package patterns

import (
	"context"

	"github.com/a-h/templ"
)

type headerActionsKey struct{}

// WithHeaderActions returns ctx carrying f, which Shell calls once, when it
// renders, to get the component shown in the page header, before the theme
// toggle (a notifications bell, say). f may return nil for none. Set it in
// middleware so every page shows the same actions; f runs only for pages that
// render a Shell.
func WithHeaderActions(ctx context.Context, f func() templ.Component) context.Context {
	return context.WithValue(ctx, headerActionsKey{}, f)
}

// HeaderActionsFrom returns the header actions for ctx, or nil.
func HeaderActionsFrom(ctx context.Context) templ.Component {
	f, ok := ctx.Value(headerActionsKey{}).(func() templ.Component)
	if !ok || f == nil {
		return nil
	}
	return f()
}
