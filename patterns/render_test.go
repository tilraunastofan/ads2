package patterns

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
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
	if !strings.Contains(html, `<h1 tabindex="-1" class="font-display text-display-sm break-words outline-none">Releases</h1>`) {
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

func TestActiveNavItem(t *testing.T) {
	groups := []NavGroup{{Items: []NavItem{
		{Label: "Home", Href: "/"},
		{Label: "Apps", Href: "/apps"},
		{Label: "Releases", Href: "/apps/releases"},
		{Label: "Help", Href: "/help", Children: []NavItem{{Label: "Guides", Href: "/help/guides"}}},
	}}}
	tests := []struct{ cur, want string }{
		{"/", "/"},
		{"/apps", "/apps"},
		{"/apps/payments", "/apps"},             // a page below an item
		{"/apps/releases/42", "/apps/releases"}, // the longest match wins over /apps
		{"/apps-old", ""},                       // a segment prefix, not a string prefix
		{"/help/guides/deploy", "/help/guides"}, // sub-items match like items
		{"/settings", ""},                       // "/" matches only itself
		{"", ""},
	}
	for _, tt := range tests {
		if got := activeHref(groups, tt.cur); got != tt.want {
			t.Errorf("activeHref(%q) = %q, want %q", tt.cur, got, tt.want)
		}
	}
}

func TestShellNavChildren(t *testing.T) {
	page := ShellProps{Title: "T", CurrentPath: "/help/guides", Nav: []NavItem{
		{Label: "Apps", Href: "/apps", Icon: templ.NopComponent},
		{Label: "Help", Href: "/help", Icon: templ.NopComponent, Children: []NavItem{
			{Label: "Guides", Href: "/help/guides"},
			{Label: "FAQ", Href: "/help/faq"},
		}},
	}}
	html := render(t, Shell(page))
	if n := strings.Count(html, `data-slot="sidebar-menu-sub-button"`); n != 2 {
		t.Fatalf("%d sub-items, want 2", n)
	}
	// The sub-list hides when the sidebar collapses to icons.
	if !strings.Contains(html, `data-slot="sidebar-menu-sub"`) || !strings.Contains(html, "group-data-[collapsible=icon]:hidden") {
		t.Error("sub-items are not in a sidebar sub-menu")
	}
	// The active sub-item and its parent are marked; nothing else is.
	var active []string
	for _, m := range regexp.MustCompile(`<a [^>]*>`).FindAllString(html, -1) {
		if regexp.MustCompile(`\sdata-active[\s>]`).MatchString(m) {
			active = append(active, regexp.MustCompile(`href="([^"]*)"`).FindStringSubmatch(m)[1])
		}
	}
	if want := []string{"/help", "/help/guides"}; !slices.Equal(active, want) {
		t.Errorf("active links %v, want %v", active, want)
	}
}

func TestShellBanner(t *testing.T) {
	page := ShellProps{Title: "T", Nav: []NavItem{{Label: "A", Href: "/a"}}}
	shell := func(b *Banner) string {
		var s strings.Builder
		ctx := context.Background()
		if b != nil {
			ctx = WithBanner(ctx, func() *Banner { return b })
		}
		if err := Shell(page).Render(ctx, &s); err != nil {
			t.Fatal(err)
		}
		return s.String()
	}
	if strings.Contains(shell(nil), `id="page-banner"`) {
		t.Error("no banner set, but one rendered")
	}
	warn := shell(&Banner{Text: "Deploys are paused.", LinkText: "Status", Href: "/status"})
	for _, want := range []string{`id="page-banner"`, "text-warning", "lucide-triangle-alert", "Deploys are paused.", `href="/status"`, ">Status</a>"} {
		if !strings.Contains(warn, want) {
			t.Errorf("warning banner lacks %s", want)
		}
	}
	bad := shell(&Banner{Destructive: true, Text: "The cluster is unreachable."})
	for _, want := range []string{"text-destructive", "lucide-circle-alert", "The cluster is unreachable."} {
		if !strings.Contains(bad, want) {
			t.Errorf("destructive banner lacks %s", want)
		}
	}
	if strings.Contains(bad[strings.Index(bad, `id="page-banner"`):], `data-slot="alert-description"`) {
		t.Error("a banner without Href should render no link line")
	}
	// The banner sits in the page column, before the page's content.
	if strings.Index(warn, `id="page-banner"`) < strings.Index(warn, `id="main"`) {
		t.Error("banner is outside #main")
	}
}

func TestShellHeaderActions(t *testing.T) {
	page := ShellProps{Title: "T", Nav: []NavItem{{Label: "A", Href: "/a"}}}
	shell := func(ctx context.Context) string {
		var s strings.Builder
		if err := Shell(page).Render(ctx, &s); err != nil {
			t.Fatal(err)
		}
		return s.String()
	}
	plain := shell(context.Background())
	if strings.Contains(plain, `id="bell"`) {
		t.Error("no header actions set, but some rendered")
	}
	// Unset: ThemeToggle stays a direct child of <header>, no wrapper.
	hdr := plain[strings.Index(plain, "<header"):strings.Index(plain, "</header>")]
	if strings.Contains(hdr, "gap-2\"><button") || strings.Contains(hdr, "<div") {
		t.Error("unset header actions added a wrapper to <header>")
	}
	calls := 0
	ctx := WithHeaderActions(context.Background(), func() templ.Component {
		calls++
		return templ.Raw(`<button id="bell">Notifications</button>`)
	})
	html := shell(ctx)
	if calls != 1 {
		t.Errorf("f called %d times in one render, want 1", calls)
	}
	bell, toggle, header := strings.Index(html, `id="bell"`), strings.Index(html, `data-theme-toggle`), strings.Index(html, "<header")
	if bell < 0 {
		t.Fatal("header actions not rendered")
	}
	if bell < header || bell > strings.Index(html, "</header>") {
		t.Error("header actions are outside <header>")
	}
	if toggle < 0 || bell > toggle {
		t.Error("header actions come after the theme toggle")
	}
	if toggle < header || toggle > strings.Index(html, "</header>") {
		t.Error("theme toggle is outside <header> when actions are set")
	}
	nilCtx := WithHeaderActions(context.Background(), func() templ.Component { return nil })
	headerOf := func(h string) string {
		// ids are random per render
		return regexp.MustCompile(`id-[A-Z0-9]+`).ReplaceAllString(h[strings.Index(h, "<header"):strings.Index(h, "</header>")], "")
	}
	if headerOf(shell(nilCtx)) != headerOf(plain) {
		t.Error("nil header actions changed the markup")
	}
}

// A wide table or code block scrolls inside its card instead of widening the
// page: the inset can shrink below its content (min-w-0 on a flex item).
func TestShellInsetCanShrink(t *testing.T) {
	html := render(t, Shell(ShellProps{Title: "T", Nav: []NavItem{{Label: "A", Href: "/a"}}}))
	i := strings.Index(html, `data-slot="sidebar-inset"`)
	if i < 0 {
		t.Fatal("no sidebar inset")
	}
	tag := html[strings.LastIndex(html[:i], "<"):]
	tag = tag[:strings.Index(tag, ">")]
	if !strings.Contains(tag, "min-w-0") {
		t.Errorf("sidebar inset lacks min-w-0: %s", tag)
	}
}

func TestShellFooter(t *testing.T) {
	page := ShellProps{
		Title:       "T",
		CurrentPath: "/about",
		Nav:         []NavItem{{Label: "Apps", Href: "/apps", Icon: templ.NopComponent}},
		Footer:      []NavItem{{Label: "About", Href: "/about", Icon: templ.NopComponent}},
		FooterNote:  "Version 1.2.3 (abc1234)",
	}
	html := render(t, Shell(page))
	_, footer, ok := strings.Cut(html, `data-slot="sidebar-footer"`)
	if !ok {
		t.Fatal("no sidebar footer")
	}
	// The footer comes after the scrolling content, so it stays at the bottom.
	if strings.Index(html, `data-slot="sidebar-content"`) > strings.Index(html, `data-slot="sidebar-footer"`) {
		t.Error("footer is not after the sidebar content")
	}
	if !strings.Contains(footer, `href="/about"`) || strings.Contains(footer, `href="/apps"`) {
		t.Error("the footer should hold the footer items and only those")
	}
	// The note is muted and hidden on the collapsed rail, where it does not fit.
	if !regexp.MustCompile(`<p class="[^"]*text-muted-foreground[^"]*group-data-\[collapsible=icon\]:hidden[^"]*">Version 1.2.3 \(abc1234\)</p>`).MatchString(footer) {
		t.Error("footer note missing, not muted, or not hidden when collapsed")
	}
	// A footer item is marked active like any other; nothing else is.
	var active []string
	for _, m := range regexp.MustCompile(`<a [^>]*>`).FindAllString(html, -1) {
		if regexp.MustCompile(`\sdata-active[\s>]`).MatchString(m) {
			active = append(active, regexp.MustCompile(`href="([^"]*)"`).FindStringSubmatch(m)[1])
		}
	}
	if want := []string{"/about"}; !slices.Equal(active, want) {
		t.Errorf("active links %v, want %v", active, want)
	}

	// Without Footer and FooterNote nothing is added, so existing apps look the same.
	if bare := render(t, Shell(ShellProps{Title: "T", Nav: page.Nav})); strings.Contains(bare, `data-slot="sidebar-footer"`) {
		t.Error("empty Footer and FooterNote should render no sidebar footer")
	}
	// The note alone renders without an empty menu.
	noteOnly := render(t, Shell(ShellProps{Title: "T", Nav: page.Nav, FooterNote: "Version dev"}))
	if _, f, _ := strings.Cut(noteOnly, `data-slot="sidebar-footer"`); !strings.Contains(f, "Version dev") || strings.Contains(f, `data-slot="sidebar-menu"`) {
		t.Error("a note without items should render the note and no menu")
	}
}
