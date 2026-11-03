package patterns

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// These tests pin what the patterns promise in docs/patterns.md, so a shadcn-templ
// upgrade that changes the markup underneath fails here instead of in an app.

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestStatusBadgeAlwaysHasIconAndLabel(t *testing.T) {
	for s, icon := range map[Status]string{
		StatusPositive:    "lucide-circle-check",
		StatusWarning:     "lucide-triangle-alert",
		StatusDestructive: "lucide-circle-alert",
		StatusNeutral:     "lucide-circle-dashed",
		Status("unknown"): "lucide-circle-dashed",
	} {
		html := render(t, StatusBadge(s, "Label "+string(s)))
		if !strings.Contains(html, icon) {
			t.Errorf("%s: no %s icon (colour must not carry the state alone)", s, icon)
		}
		if !strings.Contains(html, "Label "+string(s)) {
			t.Errorf("%s: label missing", s)
		}
	}
	if html := render(t, StatusBadge(StatusPositive, "Healthy")); !strings.Contains(html, "text-positive") {
		t.Error("positive badge is not in the positive token")
	}
	if html := render(t, StatusBadge(StatusDestructive, "Degraded")); !strings.Contains(html, `data-variant="destructive"`) {
		t.Error("destructive badge does not use the destructive variant")
	}
}

func TestPageHeader(t *testing.T) {
	html := render(t, PageHeader(PageHeaderProps{
		Crumbs: []Crumb{{Label: "Platform", Href: "/"}, {Label: "Apps", Href: "/apps"}, {Label: "Releases", Href: "/apps/releases"}},
		Title:  "Releases",
		Status: "42 applications",
	}))
	if n := strings.Count(html, "<h1"); n != 1 {
		t.Errorf("%d <h1>, want exactly one", n)
	}
	// tabindex="-1": focus-after-swap.js focuses it after a boosted navigation.
	if !strings.Contains(html, `<h1 tabindex="-1" class="font-display text-display-sm outline-none">Releases</h1>`) {
		t.Error("title is not the display style")
	}
	if !strings.Contains(html, `href="/"`) || !strings.Contains(html, `href="/apps"`) {
		t.Error("earlier crumbs are not links")
	}
	// The last crumb is the current page even when it has an Href.
	if strings.Contains(html, `href="/apps/releases"`) || !strings.Contains(html, `aria-current="page"`) {
		t.Error("last crumb should be the current page, not a link")
	}
	if !strings.Contains(html, "42 applications") {
		t.Error("status line missing")
	}

	bare := render(t, PageHeader(PageHeaderProps{Title: "Settings"}))
	if strings.Contains(bare, "breadcrumb") || strings.Contains(bare, "<p") {
		t.Error("empty Crumbs and Status should render no breadcrumb and no status line")
	}
	if !strings.Contains(bare, `<div class="flex gap-2 empty:hidden"></div>`) {
		t.Error("a header without actions should hide its (empty) actions row")
	}
}

func TestStatTrendDirectionAndColourAreIndependent(t *testing.T) {
	for _, c := range []struct {
		good, up    bool
		colour, arw string
	}{
		{true, true, "text-positive", "lucide-trending-up"},
		{true, false, "text-positive", "lucide-trending-down"},  // falling failure rate
		{false, true, "text-destructive", "lucide-trending-up"}, // rising restore time
		{false, false, "text-destructive", "lucide-trending-down"},
	} {
		html := render(t, Stat(StatProps{Label: "L", Value: "1", Delta: "+1", Good: c.good, Up: c.up}))
		if !strings.Contains(html, c.colour) || !strings.Contains(html, c.arw) {
			t.Errorf("Good=%v Up=%v: want %s and %s", c.good, c.up, c.colour, c.arw)
		}
	}
	if html := render(t, Stat(StatProps{Label: "Deployments", Value: "128"})); strings.Contains(html, `data-slot="badge"`) {
		t.Error("empty Delta should hide the trend badge")
	}
}

func TestDataCardOptionalParts(t *testing.T) {
	bare := render(t, DataCard(DataCardProps{Title: "Applications"}))
	if strings.Contains(bare, `data-slot="card-description"`) || strings.Contains(bare, `data-slot="card-action"`) {
		t.Error("empty Description and nil Action should render nothing")
	}
	full := render(t, DataCard(DataCardProps{ID: "apps", Title: "Applications", Description: "Argo CD", Action: templ.Raw(`<span id="act"></span>`)}))
	for _, want := range []string{`id="apps"`, "Argo CD", `id="act"`, "[&amp;_tr&gt;:first-child]:pl-4"} {
		if !strings.Contains(full, want) {
			t.Errorf("DataCard lacks %s", want)
		}
	}
}

// A collapsed sidebar stays collapsed across page loads: SidebarState reads the
// cookie sidebar.js writes, and Shell renders that state.
func TestShellKeepsTheSidebarCollapsed(t *testing.T) {
	page := ShellProps{Title: "T", Nav: []NavItem{{Label: "A", Href: "/a", Icon: templ.NopComponent}}}
	h := SidebarState(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = Shell(page).Render(r.Context(), w)
	}))
	get := func(cookie string) string {
		r := httptest.NewRequest("GET", "/a", nil)
		if cookie != "" {
			r.AddCookie(&http.Cookie{Name: "sidebar_state", Value: cookie})
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Body.String()
	}
	if b := get("false"); !strings.Contains(b, `data-state="collapsed"`) || !strings.Contains(b, `data-collapsible="icon"`) {
		t.Error("sidebar_state=false should render the sidebar collapsed to icons")
	}
	for _, c := range []string{"", "true"} {
		if b := get(c); strings.Contains(b, `data-state="collapsed"`) {
			t.Errorf("cookie %q should render the sidebar expanded", c)
		}
	}
}

func TestShellRendersANavItemWithoutAnIcon(t *testing.T) {
	page := ShellProps{Title: "T", Nav: []NavItem{{Label: "No icon", Href: "/a"}}}
	var b strings.Builder
	if err := Shell(page).Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), "No icon") {
		t.Error("item missing")
	}
}
