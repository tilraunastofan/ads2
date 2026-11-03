package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/tilraunastofan/ads2/patterns"
	"github.com/tilraunastofan/ads2/ui/layouts"
	"golang.org/x/net/html"
)

func get(h http.Handler, path string, hdr ...string) *httptest.ResponseRecorder {
	r := httptest.NewRequest("GET", path, nil)
	for i := 0; i+1 < len(hdr); i += 2 {
		r.Header.Set(hdr[i], hdr[i+1])
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestFontsAreCachedForAYear(t *testing.T) {
	w := get(assetHandler(), "/fonts/hanken-grotesk-latin-wght-normal.woff2")
	if w.Code != 200 {
		t.Fatalf("status %d", w.Code)
	}
	if cc := w.Header().Get("Cache-Control"); cc != "public, max-age=31536000, immutable" {
		t.Errorf("Cache-Control = %q", cc)
	}
}

func TestEmbeddedAssetsRevalidateWithETag(t *testing.T) {
	h := assetHandler()
	first := get(h, "/css/fonts.css")
	tag := first.Header().Get("ETag")
	if first.Code != 200 || tag == "" || first.Header().Get("Cache-Control") != "no-cache" {
		t.Fatalf("status %d, ETag %q, Cache-Control %q", first.Code, tag, first.Header().Get("Cache-Control"))
	}
	if again := get(h, "/css/fonts.css", "If-None-Match", tag); again.Code != http.StatusNotModified {
		t.Errorf("revalidation status %d, want 304", again.Code)
	}
}

func TestMissingAsset(t *testing.T) {
	if w := get(assetHandler(), "/nope.css"); w.Code != 404 {
		t.Errorf("status %d", w.Code)
	}
}

func TestMetricsPage(t *testing.T) {
	h := secure(newMux(newStore()))
	w := get(h, "/metrics")
	if w.Code != 200 {
		t.Fatalf("status %d", w.Code)
	}
	body := w.Body.String()
	for _, want := range []string{"Deployments per week", `id="range-label"`, `data-date-picker`, "Releases"} {
		if !strings.Contains(body, want) {
			t.Errorf("body lacks %q", want)
		}
	}
	// The CSP is style-src 'self': an inline <style> (chart.Config Color) would be blocked.
	if strings.Contains(body, "<style") {
		t.Error("page contains an inline <style>, which the CSP blocks")
	}
}

func TestMetricsFilter(t *testing.T) {
	h := newMux(newStore())
	all := get(h, "/metrics/body?from=2026-09-01&to=2026-10-12").Body.String()
	one := get(h, "/metrics/body?from=2026-09-01&to=2026-10-12&app=payments-api").Body.String()
	if !strings.Contains(all, `<hx-partial hx-target="#range-label"`) {
		t.Error("filter response lacks the range label partial")
	}
	if !strings.Contains(all, "1 Sep") || !strings.Contains(all, "12 Oct") {
		t.Error("range label does not reflect from and to")
	}
	if rows := strings.Count(one, "<tr"); rows != 2 || strings.Count(all, "<tr") <= rows {
		t.Errorf("app filter: %d rows filtered (want header + 1), %d unfiltered", rows, strings.Count(all, "<tr"))
	}
	if !strings.Contains(one, "payments-api") {
		t.Error("filtered response lacks the app")
	}
}

func TestReleaseDetails(t *testing.T) {
	h := newMux(newStore())
	if w := get(h, "/metrics/release?app=payments-api&version=1.42.0"); w.Code != 200 || !strings.Contains(w.Body.String(), "registry.internal/arion/payments-api:1.42.0") {
		t.Errorf("known release: status %d", w.Code)
	}
	if w := get(h, "/metrics/release?app=nope&version=9"); w.Code != 404 {
		t.Errorf("unknown release: status %d, want 404", w.Code)
	}
}

func post(h http.Handler, path, form string) *httptest.ResponseRecorder {
	r := httptest.NewRequest("POST", path, strings.NewReader(form))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestSettingsApproval(t *testing.T) {
	h := newMux(newStore())
	if w := get(h, "/settings"); w.Code != 200 || !strings.Contains(w.Body.String(), "Approval code") {
		t.Fatalf("GET /settings: status %d", w.Code)
	}
	bad := post(h, "/settings/approve", "code=111111")
	// The input and each of the six slots: the slots' red border does not follow the input.
	if n := strings.Count(bad.Body.String(), `aria-invalid="true"`); bad.Code != http.StatusUnprocessableEntity || n != 7 {
		t.Errorf("wrong code: status %d, %d aria-invalid=\"true\", want 7", bad.Code, n)
	}
	if strings.Contains(get(h, "/settings").Body.String(), `aria-invalid="true"`) {
		t.Error("fresh form is marked invalid")
	}
	ok := post(h, "/settings/approve", "code=123456")
	if ok.Code != 200 || !strings.Contains(ok.Body.String(), "Freeze lifted") || !strings.Contains(ok.Body.String(), `hx-target="#toasts"`) {
		t.Errorf("right code: status %d, body %q", ok.Code, ok.Body.String())
	}
}

func TestAppsCSV(t *testing.T) {
	h := newMux(newStore())
	w := get(h, "/apps.csv")
	if w.Code != 200 || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/csv") {
		t.Fatalf("status %d, Content-Type %q", w.Code, w.Header().Get("Content-Type"))
	}
	if !strings.Contains(w.Header().Get("Content-Disposition"), "attachment") {
		t.Error("CSV is not served as a download")
	}
	body := w.Body.String()
	if !strings.HasPrefix(body, "application,environment,version,health\n") || !strings.Contains(body, "payments-api,Production,1.42.0,Healthy\n") {
		t.Errorf("body %q", body)
	}
	// The Export button must be a download link, or htmx's boost swaps the CSV into the page.
	if page := get(h, "/apps").Body.String(); !strings.Contains(page, `download="applications.csv"`) {
		t.Error("Export link lacks the download attribute")
	}
}

func TestPartialEscapesAttributes(t *testing.T) {
	var b strings.Builder
	if err := partial(`#x" onclick="y`, "outerHTML", templ.Raw("")).Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(b.String(), `" onclick="`) {
		t.Errorf("target broke out of the attribute: %s", b.String())
	}
}

func TestShellAccessibility(t *testing.T) {
	body := get(newMux(newStore()), "/apps").Body.String()
	skip := strings.Index(body, `href="#main"`)
	if skip < 0 || skip > strings.Index(body, `data-sidebar="trigger"`) {
		t.Error("skip link missing or not before the navigation")
	}
	// One <main> landmark: the sidebar inset. The skip link's target is inside it.
	if n := strings.Count(body, "<main"); n != 1 {
		t.Errorf("%d <main> landmarks, want 1", n)
	}
	for _, want := range []string{`<div id="main" tabindex="-1"`, `aria-label="Toggle sidebar"`, `aria-label="Dark theme"`, `aria-pressed="false"`} {
		if !strings.Contains(body, want) {
			t.Errorf("page lacks %s", want)
		}
	}
}

func gunzip(t *testing.T, b []byte) string {
	t.Helper()
	zr, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(zr)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func TestAssetsArePrecompressed(t *testing.T) {
	h := compress(assetHandler())
	plain := get(h, "/css/fonts.css")
	gz := get(h, "/css/fonts.css", "Accept-Encoding", "br, gzip")
	if gz.Header().Get("Content-Encoding") != "gzip" || !strings.Contains(gz.Header().Get("Vary"), "Accept-Encoding") {
		t.Fatalf("Content-Encoding %q, Vary %q", gz.Header().Get("Content-Encoding"), gz.Header().Get("Vary"))
	}
	if got := gunzip(t, gz.Body.Bytes()); got != plain.Body.String() {
		t.Error("gzipped asset differs from the plain one")
	}
	if !strings.HasPrefix(gz.Header().Get("Content-Type"), "text/css") {
		t.Errorf("Content-Type %q", gz.Header().Get("Content-Type"))
	}
	tag := gz.Header().Get("ETag")
	if tag == plain.Header().Get("ETag") {
		t.Error("gzip and identity share an ETag")
	}
	if again := get(h, "/css/fonts.css", "Accept-Encoding", "gzip", "If-None-Match", tag); again.Code != http.StatusNotModified {
		t.Errorf("revalidation of gzipped asset: %d, want 304", again.Code)
	}
	if w := get(h, "/css/fonts.css", "Accept-Encoding", "gzip;q=0"); w.Header().Get("Content-Encoding") != "" {
		t.Error("gzip;q=0 still got gzip")
	}
	// woff2 is already compressed.
	if w := get(h, "/fonts/hanken-grotesk-latin-wght-normal.woff2", "Accept-Encoding", "gzip"); w.Header().Get("Content-Encoding") != "" {
		t.Error("font was gzipped")
	}
}

func TestPagesAreCompressed(t *testing.T) {
	h := secure(compress(newMux(newStore())))
	w := get(h, "/apps", "Accept-Encoding", "gzip")
	if w.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("Content-Encoding %q", w.Header().Get("Content-Encoding"))
	}
	if body := gunzip(t, w.Body.Bytes()); !strings.Contains(body, "Sync and health") {
		t.Error("decompressed page lacks its title")
	}
	if w := get(h, "/apps"); w.Header().Get("Content-Encoding") != "" || !strings.Contains(w.Body.String(), "Sync and health") {
		t.Error("client without gzip got a compressed or broken page")
	}
}

func TestNotFoundPage(t *testing.T) {
	h := newMux(newStore())
	w := get(h, "/metrics/nope")
	body := w.Body.String()
	if w.Code != http.StatusNotFound {
		t.Fatalf("status %d", w.Code)
	}
	for _, want := range []string{"Page not found", `<code class="font-mono text-code break-all text-foreground">/metrics/nope</code>`, `href="/metrics"`, "Go to Metrics", `data-sidebar="trigger"`} {
		if !strings.Contains(body, want) {
			t.Errorf("404 page lacks %s", want)
		}
	}
	if !strings.Contains(get(h, "/nope").Body.String(), "Go to Applications") {
		t.Error("unknown section should fall back to Applications")
	}
	// A broken hx-get must not swap a whole page into its target.
	if p := get(h, "/nope", "HX-Request", "true", "HX-Request-Type", "partial"); p.Code != 404 || strings.Contains(p.Body.String(), "<html") {
		t.Errorf("partial 404: status %d, body %q", p.Code, p.Body.String())
	}
	if p := get(h, "/nope", "HX-Request", "true", "HX-Request-Type", "full"); !strings.Contains(p.Body.String(), "Page not found") {
		t.Error("boosted navigation should get the 404 page")
	}
}

// The sheet's lazy body must describe the release in the row it opens from, for
// the selected range, not for the default one.
func TestReleaseDetailsFollowTheSelectedRange(t *testing.T) {
	h := newMux(newStore())
	body := get(h, "/metrics/body?from=2026-09-08&to=2026-09-20&app=fx-rates").Body.String()
	if !strings.Contains(body, "18 Sep 2026, 11:00") {
		t.Fatalf("table lacks the fx-rates release time for the range")
	}
	const want = `hx-get="/metrics/release?app=fx-rates&amp;from=2026-09-08&amp;to=2026-09-20&amp;version=3.1.0"`
	if !strings.Contains(body, want) {
		t.Fatalf("sheet URL lacks the range; want %s", want)
	}
	details := get(h, "/metrics/release?app=fx-rates&from=2026-09-08&to=2026-09-20&version=3.1.0").Body.String()
	if !strings.Contains(details, "18 Sep 2026, 11:00") {
		t.Errorf("details show another release time than the table")
	}
}

func TestRecovererAnswersAPanic(t *testing.T) {
	boom := recoverer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") }))
	frag := get(boom, "/x", "HX-Request", "true", "HX-Request-Type", "partial")
	if frag.Code != 500 || !strings.Contains(frag.Body.String(), `<hx-partial hx-target="#toasts"`) || !strings.Contains(frag.Body.String(), "Something went wrong") {
		t.Errorf("htmx request: status %d, body %q", frag.Code, frag.Body.String())
	}
	page := get(boom, "/x")
	if page.Code != 500 || !strings.Contains(page.Body.String(), "<h1") || strings.Contains(page.Body.String(), "boom") {
		t.Errorf("page request: status %d; the panic value must not reach the client", page.Code)
	}
	if ok := get(recoverer(newMux(newStore())), "/apps"); ok.Code != 200 {
		t.Errorf("recoverer broke a normal request: %d", ok.Code)
	}
}

func TestShellHasPageAlertsAndHelpers(t *testing.T) {
	body := get(newMux(newStore()), "/apps").Body.String()
	for _, want := range []string{`<div id="page-alerts" role="status" aria-live="polite" class="empty:hidden"></div>`,
		`src="/assets/js/request-errors.js"`, `src="/assets/js/toast-guard.js"`, `src="/assets/js/filter-form.js"`} {
		if !strings.Contains(body, want) {
			t.Errorf("page lacks %s", want)
		}
	}
	if strings.Contains(body, "Target for errors") {
		t.Error("a templ comment leaked into the page")
	}
}

// primaries collects the primary buttons of a page by the view they are seen in:
// "page", the page plus one tab panel ("tab:<n>"), or an overlay ("overlay:<n>").
// A primary is a <button> or <a> with the exact class bg-primary (badges, checked
// controls and calendar days use prefixed or other slots).
func primaries(t *testing.T, body string) map[string][]string {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	out := map[string][]string{}
	n := 0
	var text func(*html.Node) string
	text = func(x *html.Node) string {
		if x.Type == html.TextNode {
			return strings.TrimSpace(x.Data)
		}
		var b strings.Builder
		for c := x.FirstChild; c != nil; c = c.NextSibling {
			b.WriteString(text(c))
		}
		return b.String()
	}
	var walk func(x *html.Node, view string)
	walk = func(x *html.Node, view string) {
		if x.Type == html.ElementNode {
			slot, cls := "", ""
			for _, a := range x.Attr {
				switch a.Key {
				case "data-slot":
					slot = a.Val
				case "class":
					cls = a.Val
				}
			}
			switch slot {
			case "dialog-content", "alert-dialog-content", "sheet-content", "popover-content", "dropdown-menu-content":
				n++
				view = "overlay:" + strconv.Itoa(n)
			case "tabs-content":
				if view == "page" {
					n++
					view = "tab:" + strconv.Itoa(n)
				}
			}
			if (x.Data == "button" || x.Data == "a") && strings.Contains(" "+cls+" ", " bg-primary ") {
				out[view] = append(out[view], text(x))
			}
		}
		for c := x.FirstChild; c != nil; c = c.NextSibling {
			walk(c, view)
		}
	}
	walk(doc, "page")
	return out
}

func TestOnePrimaryButtonPerView(t *testing.T) {
	h := newMux(newStore())
	for _, path := range []string{"/apps", "/metrics", "/settings", "/nope"} {
		views := primaries(t, get(h, path).Body.String())
		page := views["page"]
		if len(page) > 1 {
			t.Errorf("%s: %d primary buttons on the page: %q", path, len(page), page)
		}
		for v, b := range views {
			switch {
			case strings.HasPrefix(v, "tab:") && len(page)+len(b) > 1:
				t.Errorf("%s: %s shows %q next to the page's %q", path, v, b, page)
			case strings.HasPrefix(v, "overlay:") && len(b) > 1:
				t.Errorf("%s: %s has %d primary buttons: %q", path, v, len(b), b)
			}
		}
	}
}

func TestPrimaryCounterSeesViolations(t *testing.T) {
	const two = `<main><button class="x bg-primary y">A</button><a class="bg-primary">B</a>` +
		`<span data-slot="badge" class="bg-primary">badge</span><button class="data-checked:bg-primary">c</button>` +
		`<div data-slot="tabs-content"><button class="bg-primary">T</button></div>` +
		`<div data-slot="dialog-content"><button class="bg-primary">D</button></div></main>`
	v := primaries(t, two)
	if len(v["page"]) != 2 || len(v["tab:1"]) != 1 || len(v["overlay:2"]) != 1 {
		t.Errorf("got %v", v)
	}
}

func TestIcelandicReport(t *testing.T) {
	w := get(newMux(newStore()), "/uppgjor")
	body := w.Body.String()
	if w.Code != 200 {
		t.Fatalf("status %d", w.Code)
	}
	for _, want := range []string{
		`<html lang="is">`, "Uppgjör 3F 2026", "Birt 23. október 2026", "Fara beint í efni", `aria-label="Dökkt þema"`,
		"7.412 m.kr.", "+4,9%", "14,2%", "+0,9 pst.", "41,6%", "−0,9 pst.", // KPIs
		"16.711", "15.947", "9.252", "−6.947", "−2,5%", // computed subtotals, a cost line
		"Í milljónum króna (m.kr.)", "Tölurnar eru sýnigögn",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("report lacks %q", want)
		}
	}
	for _, en := range []string{"Skip to content", "Dark theme", "Toggle sidebar", "Applications", "Workloads"} {
		if strings.Contains(body, en) {
			t.Errorf("Icelandic page contains English %q", en)
		}
	}
	// A cost that grew is a loss: minus and destructive; the falling cost ratio is good.
	if !strings.Contains(body, `text-destructive">−2,5%`) {
		t.Error("growing cost is not shown as a loss")
	}
	if !regexp.MustCompile(`(?s)text-positive"><svg[^>]*lucide-trending-down.*?−0,9 pst\.`).MatchString(body) {
		t.Error("falling cost ratio should be good (positive) with a down arrow")
	}
}

// One name per thing: the sidebar label, the breadcrumb and the page title agree,
// and no sidebar item is a second way to the same page.
func TestNavLabelsMatchPageTitles(t *testing.T) {
	h := newMux(newStore())
	titles := map[string]string{"/apps": "Applications", "/metrics": "Metrics", "/settings": "Deployment freeze"}
	for _, it := range layouts.Nav() {
		want, ok := titles[it.Href]
		if !ok {
			continue // the report page has its own layout
		}
		if it.Label != want {
			t.Errorf("sidebar says %q for %s, page is %q", it.Label, it.Href, want)
		}
		body := get(h, it.Href).Body.String()
		if !strings.Contains(body, `<h1 class="font-display text-display-sm">`+want+`</h1>`) {
			t.Errorf("%s: title is not %q", it.Href, want)
		}
	}
	seen := map[string]bool{}
	for _, it := range layouts.Nav() {
		if seen[it.Href] || it.Href == "/" {
			t.Errorf("sidebar item %q duplicates another destination (%s)", it.Label, it.Href)
		}
		seen[it.Href] = true
	}
}

// Delete and Scale live in the row menu and open the dialog for that row.
func TestRowMenuOpensTheRowsDialogs(t *testing.T) {
	body := get(newMux(newStore()), "/apps").Body.String()
	if strings.Contains(body, `aria-label="Delete payments-api"`) {
		t.Error("a delete button is still shown outside the row menu")
	}
	for _, id := range []string{"delete-payments-api", "scale-payments-api"} {
		if !strings.Contains(body, `data-templ-controls="`+id+`"`) || !strings.Contains(body, `id="`+id+`"`) {
			t.Errorf("menu item or dialog for %s missing", id)
		}
	}
	if !strings.Contains(body, "Scale payments-api") {
		t.Error("scale dialog does not name its application")
	}
}

func TestCreateApplicationAddsARow(t *testing.T) {
	h := newMux(newStore())
	ok := post(h, "/apps", "name=ledger-sync&replicas=3&namespace=payments")
	b := ok.Body.String()
	if ok.Code != 200 || !strings.Contains(b, `<hx-partial hx-target="#app-rows"`) || !strings.Contains(b, `id="app-ledger-sync"`) {
		t.Fatalf("create: status %d, row partial missing", ok.Code)
	}
	if dup := post(h, "/apps", "name=ledger-sync&replicas=3"); dup.Code != http.StatusUnprocessableEntity || !strings.Contains(dup.Body.String(), "already exists") {
		t.Errorf("duplicate name: status %d", dup.Code)
	}
	if !strings.Contains(get(h, "/apps").Body.String(), `id="app-ledger-sync"`) {
		t.Error("new application missing after reload")
	}
}

func TestScaleToastNamesTheApplication(t *testing.T) {
	w := post(newMux(newStore()), "/scale", "app=fx-rates&replicas=4")
	if !strings.Contains(w.Body.String(), "Scale fx-rates to 4 replicas.") {
		t.Errorf("toast: %q", w.Body.String())
	}
}

// The Applications KPIs and the Metrics chart read the same week data.
func TestApplicationsKPIsMatchMetrics(t *testing.T) {
	now := time.Now()
	_, wk := now.ISOWeek()
	d, _ := week(wk)
	if got := appsKPIs(now)[0].Value; got != patterns.Int(int64(d)) {
		t.Errorf("Deploys this week = %s, week data says %d", got, d)
	}
	body := get(newMux(newStore()), "/metrics").Body.String()
	if !strings.Contains(body, "W"+strconv.Itoa(wk)) {
		t.Errorf("metrics chart lacks this week W%d", wk)
	}
}

func TestFreezeCopyKeepsOneName(t *testing.T) {
	h := newMux(newStore())
	if !strings.Contains(get(h, "/settings").Body.String(), "Lift freeze") {
		t.Error("button is not Lift freeze")
	}
	if !strings.Contains(post(h, "/settings/approve", "code=123456").Body.String(), "Freeze lifted") {
		t.Error("success does not say Freeze lifted")
	}
}

// The deploy chart has one lead: the latest week in full chart-1, earlier weeks in
// chart-1-quiet, rollbacks in neutral chart-4. The chart is drawn in the browser
// from the JSON model in the response, so the test reads that.
func TestDeployChartHighlightsTheLatestWeek(t *testing.T) {
	// 7 Sep to 4 Oct 2026 is ISO weeks 37 to 40.
	body := get(newMux(newStore()), "/metrics/body?from=2026-09-07&to=2026-10-04").Body.String()
	const cells = `"cells":["var(--chart-1-quiet)","var(--chart-1-quiet)","var(--chart-1-quiet)","var(--chart-1)"]`
	if !strings.Contains(body, cells) {
		t.Errorf("deployments bar cells: want %s", cells)
	}
	if !strings.Contains(body, `"key":"rollbacks","label":"Rollbacks","color":"var(--chart-4)"`) || strings.Contains(body, "var(--chart-2)") {
		t.Error("rollbacks should be chart-4 (slate), and gold chart-2 is no longer used")
	}
}

func TestNewAppFieldErrorsAreAnnounced(t *testing.T) {
	h := newMux(newStore())
	b := post(h, "/apps", "name=Bad_Name&replicas=1").Body.String()
	// A bare aria-invalid reads as false; it must be the string "true", on both inputs.
	if n := strings.Count(b, `aria-invalid="true"`); n != 2 {
		t.Errorf("%d aria-invalid=\"true\", want 2", n)
	}
	for _, id := range []string{"name", "replicas"} {
		if !strings.Contains(b, `aria-describedby="`+id+`-hint"`) || !strings.Contains(b, `id="`+id+`-hint"`) {
			t.Errorf("%s: error not tied to the input by aria-describedby", id)
		}
	}
	fresh := get(h, "/apps").Body.String()
	if strings.Contains(fresh, `aria-invalid="true"`) || strings.Contains(fresh, "aria-invalid ") || strings.Contains(fresh, "aria-invalid>") {
		t.Error("fresh form is marked invalid")
	}
	if !strings.Contains(fresh, `aria-describedby="name-hint"`) {
		t.Error("description not tied to the name input")
	}
}

func TestAppNamesAreDNSLabels(t *testing.T) {
	h := newMux(newStore())
	// Each would break the row's /apps/{name} URLs, or run as a spreadsheet formula
	// in /apps.csv.
	for _, name := range []string{"a/b", "a?b", "a#b", "=1+1", "+1", "-app", "app-", "a.b", "Ab", strings.Repeat("a", 64), ""} {
		if w := post(h, "/apps", "replicas=3&name="+url.QueryEscape(name)); w.Code != http.StatusUnprocessableEntity {
			t.Errorf("name %q: status %d, want 422", name, w.Code)
		}
	}
	for _, name := range []string{"a", "ledger-sync", "app2", strings.Repeat("a", 63)} {
		if w := post(h, "/apps", "replicas=3&name="+name); w.Code != 200 {
			t.Errorf("name %q: status %d, want 200", name, w.Code)
		}
	}
}

func TestCSVCellsNeverRunAsFormulas(t *testing.T) {
	for in, want := range map[string]string{"=1+1": "'=1+1", "+1": "'+1", "-1": "'-1", "@x": "'@x", "\tx": "'\tx", "ledger": "ledger", "": ""} {
		if got := csvCell(in); got != want {
			t.Errorf("csvCell(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMetricsRangeIsBounded(t *testing.T) {
	h := newMux(newStore())
	// A crafted range of 10,000 years would be one bar per week of it.
	w := get(h, "/metrics/body?from=0001-01-01&to=9999-12-31")
	if w.Code != 200 {
		t.Fatalf("status %d", w.Code)
	}
	if n := strings.Count(w.Body.String(), `"week":`); n > 54 {
		t.Errorf("%d weeks rendered, want at most 54", n)
	}
	from, to := boundRange(time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC), time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC))
	if !from.Before(to) || from.Month() != 9 || to.Month() != 10 {
		t.Errorf("reversed range: got %s to %s, want it put the right way round", from, to)
	}
}
