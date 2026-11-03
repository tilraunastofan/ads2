package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"
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
	if w := get(h, "/metrics/release?app=payments-api&version=1.40.0"); w.Code != 200 || !strings.Contains(w.Body.String(), "registry.internal/arion/payments-api:1.40.0") {
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
	for _, want := range []string{`<main id="main" tabindex="-1"`, `aria-label="Toggle sidebar"`, `aria-label="Dark theme"`, `aria-pressed="false"`} {
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
	if body := gunzip(t, w.Body.Bytes()); !strings.Contains(body, "Releases") {
		t.Error("decompressed page lacks its title")
	}
	if w := get(h, "/apps"); w.Header().Get("Content-Encoding") != "" || !strings.Contains(w.Body.String(), "Releases") {
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
