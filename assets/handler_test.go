package assets

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func serve(h http.Handler, path string, hdr ...string) *httptest.ResponseRecorder {
	r := httptest.NewRequest("GET", path, nil)
	for i := 0; i+1 < len(hdr); i += 2 {
		r.Header.Set(hdr[i], hdr[i+1])
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestOverlayPutsTheAppOnTop(t *testing.T) {
	app := fstest.MapFS{
		"css/output.css": {Data: []byte("/* app build */")},
		"css/fonts.css":  {Data: []byte("/* app override */")},
	}
	h := Handler(Overlay(app))
	if b := serve(h, "/css/output.css").Body.String(); b != "/* app build */" {
		t.Errorf("app file: got %q", b)
	}
	if b := serve(h, "/css/fonts.css").Body.String(); b != "/* app override */" {
		t.Errorf("an app file of the same name must win: got %q", b)
	}
	// Everything else comes from the design system.
	for _, p := range []string{"/css/arion-theme.css", "/js/htmx.min.js", "/js/theme.js", "/img/arion-mark.svg"} {
		if w := serve(h, p); w.Code != 200 {
			t.Errorf("%s: status %d, want the design system's file", p, w.Code)
		}
	}
	if w := serve(h, "/css/nope.css"); w.Code != 404 {
		t.Errorf("missing file: status %d", w.Code)
	}
}

func TestEmbeddedAssetsAreComplete(t *testing.T) {
	// What patterns.Base links to, so an app on Handler(Overlay(app)) only adds css/output.css.
	h := Handler(Assets)
	for _, p := range []string{
		"/css/fonts.css", "/css/ads2.css", "/js/theme.js", "/js/htmx.min.js", "/js/toast-guard.js",
		"/js/filter-form.js", "/js/request-errors.js", "/js/dialog-form.js", "/js/datepicker.js",
		"/js/focus-after-swap.js", "/img/arion-mark-favicon.svg",
		"/fonts/hanken-grotesk-latin-wght-normal.woff2", "/fonts/source-serif-4-latin-wght-normal.woff2",
	} {
		if w := serve(h, p); w.Code != 200 {
			t.Errorf("%s: status %d", p, w.Code)
		}
	}
}

func TestCacheAndGzip(t *testing.T) {
	h := Handler(Assets)
	font := serve(h, "/fonts/hanken-grotesk-latin-wght-normal.woff2")
	if cc := font.Header().Get("Cache-Control"); !strings.Contains(cc, "immutable") {
		t.Errorf("font Cache-Control %q", cc)
	}
	css := serve(h, "/css/arion-theme.css", "Accept-Encoding", "gzip")
	if css.Header().Get("Content-Encoding") != "gzip" || css.Header().Get("Cache-Control") != "no-cache" {
		t.Errorf("theme: encoding %q, cache %q", css.Header().Get("Content-Encoding"), css.Header().Get("Cache-Control"))
	}
	again := serve(h, "/css/arion-theme.css", "Accept-Encoding", "gzip", "If-None-Match", css.Header().Get("ETag"))
	if again.Code != http.StatusNotModified {
		t.Errorf("revalidation: status %d, want 304", again.Code)
	}
}

func TestAcceptsGzip(t *testing.T) {
	for hdr, want := range map[string]bool{
		"gzip": true, "br, gzip": true, "gzip;q=0.5": true, "GZIP; q=1": true,
		"gzip;q=0": false, "gzip;q=0.0": false, "gzip; q=0.000": false, "br": false, "": false, "gzip;q=x": false,
	} {
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Accept-Encoding", hdr)
		if got := AcceptsGzip(r); got != want {
			t.Errorf("AcceptsGzip(%q) = %v, want %v", hdr, got, want)
		}
	}
}
