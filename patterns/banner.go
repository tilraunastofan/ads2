package patterns

import "context"

// Banner is a page-wide notice that Shell shows above the page's content, such as
// "Deploys are paused: the cluster API is unreachable." It is for a state of the
// whole tool, not of one page; page errors go to #page-alerts.
type Banner struct {
	Destructive bool   // a failure; otherwise a warning
	Text        string // the notice, one sentence
	LinkText    string // optional link to the details; defaults to Href
	Href        string
}

type bannerKey struct{}

// WithBanner returns ctx carrying f, which Shell calls once, when it renders, to
// get the banner (nil for none). Set it in middleware so every page shows the same
// notice; f runs only for pages that render a Shell.
func WithBanner(ctx context.Context, f func() *Banner) context.Context {
	return context.WithValue(ctx, bannerKey{}, f)
}

// BannerFrom returns the banner for ctx, or nil.
func BannerFrom(ctx context.Context) *Banner {
	f, ok := ctx.Value(bannerKey{}).(func() *Banner)
	if !ok || f == nil {
		return nil
	}
	return f()
}
