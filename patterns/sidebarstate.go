package patterns

import (
	"context"
	"net/http"
)

// sidebar.js remembers the sidebar's state in this cookie when the trigger is
// clicked: "false" while collapsed, "true" while expanded (shadcn's sidebar_state).
const sidebarCookie = "sidebar_state"

type sidebarCollapsedKey struct{}

// SidebarState reads the sidebar cookie into the request context, so Shell renders
// the sidebar the way the user left it. Without it every page load, and so every
// boosted navigation (which swaps in a freshly rendered body), renders the sidebar
// expanded again. shadcn's React layout does the same with defaultOpen.
//
//	http.ListenAndServe(addr, patterns.SidebarState(mux))
func SidebarState(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(sidebarCookie); err == nil && c.Value == "false" {
			r = r.WithContext(WithSidebarCollapsed(r.Context(), true))
		}
		next.ServeHTTP(w, r)
	})
}

// WithSidebarCollapsed sets the sidebar's initial state for Shell, for code that
// renders outside SidebarState (tests, a stored preference).
func WithSidebarCollapsed(ctx context.Context, collapsed bool) context.Context {
	return context.WithValue(ctx, sidebarCollapsedKey{}, collapsed)
}

func sidebarCollapsed(ctx context.Context) bool {
	c, _ := ctx.Value(sidebarCollapsedKey{}).(bool)
	return c
}
