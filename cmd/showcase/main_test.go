package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
