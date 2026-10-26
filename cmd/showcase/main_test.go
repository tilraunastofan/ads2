package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
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
