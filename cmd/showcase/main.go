// Command showcase serves the Arion Design System 2 reference app: shadcn-templ
// components themed with assets/css/arion-theme.css, driven by htmx 4.
//
//	task dev          # templ + tailwind + script watchers, http://localhost:7331
//	PORT=8090 go run ./cmd/showcase
package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"html"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/a-h/templ"

	"github.com/tilraunastofan/ads2/assets"
	"github.com/tilraunastofan/ads2/components/toast"
	"github.com/tilraunastofan/ads2/patterns"
	"github.com/tilraunastofan/ads2/ui/pages"
)

type store struct {
	mu   sync.Mutex
	apps []pages.App
	tick int
}

func newStore() *store {
	return &store{apps: []pages.App{
		{Name: "payments-api", Env: "Production", Version: "1.42.0", Health: "Healthy"},
		{Name: "lending-web", Env: "Staging", Version: "0.19.3", Health: "Progressing"},
		{Name: "fx-rates", Env: "Production", Version: "3.1.0", Health: "Degraded"},
		{Name: "cards-ledger", Env: "Production", Version: "2.7.4", Health: "Healthy"},
		{Name: "kyc-worker", Env: "Development", Version: "0.4.0", Health: "Healthy"},
	}}
}

func (s *store) search(q string) []pages.App {
	s.mu.Lock()
	defer s.mu.Unlock()
	q = strings.ToLower(strings.TrimSpace(q))
	out := make([]pages.App, 0, len(s.apps))
	for _, a := range s.apps {
		if q == "" || strings.Contains(a.Name, q) {
			out = append(out, a)
		}
	}
	return out
}

func (s *store) remove(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, a := range s.apps {
		if a.Name == name {
			s.apps = append(s.apps[:i], s.apps[i+1:]...)
			return
		}
	}
}

func render(w http.ResponseWriter, r *http.Request, status int, c templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := c.Render(r.Context(), w); err != nil {
		log.Printf("render: %v", err)
	}
}

// toastPartial wraps a toast marker in an htmx 4 <hx-partial> aimed at #toasts,
// so any response (even hx-swap="none") can raise a toast.
func toastPartial(title, desc string, t toast.Type) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<hx-partial hx-target="#toasts" hx-swap="beforeend">`); err != nil {
			return err
		}
		if err := pages.Toast(title, desc, t).Render(ctx, w); err != nil {
			return err
		}
		_, err := io.WriteString(w, `</hx-partial>`)
		return err
	})
}

// partial wraps a component in an htmx 4 <hx-partial> for a second target. The
// attribute values are escaped, so a target can't break out of the tag.
func partial(target, swap string, c templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<hx-partial hx-target="`+html.EscapeString(target)+`" hx-swap="`+html.EscapeString(swap)+`">`); err != nil {
			return err
		}
		if err := c.Render(ctx, w); err != nil {
			return err
		}
		_, err := io.WriteString(w, `</hx-partial>`)
		return err
	})
}

// metrics builds the Metrics view model from the query (from, to, app). The data
// is generated deterministically so the showcase needs no backend.
func (s *store) metrics(r *http.Request) pages.MetricsVM {
	q := r.URL.Query()
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	to, err := patterns.ParseISODate(q.Get("to"), time.Local)
	if err != nil {
		to = today
	}
	from, err := patterns.ParseISODate(q.Get("from"), time.Local)
	if err != nil || from.After(to) {
		from = to.AddDate(0, 0, -41) // six weeks
	}
	app := q.Get("app")

	s.mu.Lock()
	apps := make([]string, 0, len(s.apps))
	for _, a := range s.apps {
		apps = append(apps, a.Name)
	}
	s.mu.Unlock()

	var weeks []pages.WeekStat
	deploys, rollbacks := 0, 0
	for d := from; !d.After(to); d = d.AddDate(0, 0, 7) {
		_, wk := d.ISOWeek()
		dep := 18 + (wk*7)%17
		rb := 1 + (wk*3)%5
		if app != "" {
			dep, rb = max(dep/5, 1), rb/3
		}
		weeks = append(weeks, pages.WeekStat{Week: "W" + strconv.Itoa(wk), Deploys: dep, Rollbacks: rb})
		deploys += dep
		rollbacks += rb
	}

	var releases []pages.Release
	health := []string{"Healthy", "Healthy", "Progressing", "Healthy", "Degraded"}
	for i, a := range apps {
		if app != "" && a != app {
			continue
		}
		releases = append(releases, pages.Release{
			App:     a,
			Version: "1." + strconv.Itoa(40+i) + ".0",
			Env:     "Production",
			At:      to.Add(time.Duration(9+i)*time.Hour).AddDate(0, 0, -i),
			By:      "Kargo · stage prod",
			Health:  health[i%len(health)],
		})
	}

	cfr := 0.0
	if deploys > 0 {
		cfr = 100 * float64(rollbacks) / float64(deploys)
	}
	return pages.MetricsVM{
		From: from, To: to, App: app, Apps: apps, Weeks: weeks, Releases: releases,
		KPIs: []patterns.StatProps{
			{Label: "Deployments", Value: patterns.Int(int64(deploys))},
			{Label: "Rollbacks", Value: patterns.Int(int64(rollbacks))},
			{Label: "Change failure rate", Value: patterns.Pct(cfr, 1)},
		},
	}
}

// assetHandler serves the embedded assets, or ./assets when SHADCN_TEMPL_DEV=true
// so the watchers' output shows up without a rebuild.
//
// Fonts and the hashed component bundle never change under their name, so they are
// cached for a year, in dev too (re-downloading fonts on every refresh makes text
// flash from the fallback font). Everything else revalidates; embedded files carry
// no modification time, so they get a content-hash ETag and a refresh costs a 304.
func assetHandler() http.Handler {
	dev := os.Getenv("SHADCN_TEMPL_DEV") == "true"
	var disk http.Handler
	if dev {
		disk = http.FileServer(http.Dir("assets"))
	}
	var cache sync.Map // asset name -> *cachedAsset
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/")
		immutable := strings.HasPrefix(name, "fonts/") || strings.HasPrefix(name, "js/shadcn-templ-")
		switch {
		case immutable:
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		case dev:
			w.Header().Set("Cache-Control", "no-store")
		default:
			w.Header().Set("Cache-Control", "no-cache")
		}
		if dev {
			disk.ServeHTTP(w, r)
			return
		}
		body, err := fs.ReadFile(assets.Assets, name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		a, ok := cache.Load(name)
		if !ok {
			a, _ = cache.LoadOrStore(name, newCachedAsset(name, body))
		}
		ca := a.(*cachedAsset)
		etag := ca.etag
		if ca.gz != nil {
			w.Header().Add("Vary", "Accept-Encoding")
			if acceptsGzip(r) {
				// The gzipped bytes are their own representation: own ETag, and
				// ServeContent applies ranges and If-None-Match to them.
				body, etag = ca.gz, ca.gzETag
				w.Header().Set("Content-Encoding", "gzip")
			}
		}
		w.Header().Set("ETag", etag)
		http.ServeContent(w, r, path.Base(name), time.Time{}, bytes.NewReader(body))
	})
}

// cachedAsset is an embedded file's ETag and, for text types, its gzipped bytes,
// computed once per process.
type cachedAsset struct {
	etag   string
	gz     []byte
	gzETag string
}

func newCachedAsset(name string, body []byte) *cachedAsset {
	sum := sha256.Sum256(body)
	h := hex.EncodeToString(sum[:8])
	a := &cachedAsset{etag: `"` + h + `"`}
	switch path.Ext(name) {
	case ".css", ".js", ".svg", ".json", ".txt":
		var b bytes.Buffer
		zw, _ := gzip.NewWriterLevel(&b, gzip.BestCompression)
		_, _ = zw.Write(body)
		if zw.Close() == nil && b.Len() < len(body) {
			a.gz, a.gzETag = b.Bytes(), `"`+h+`-gz"`
		}
	}
	return a
}

func acceptsGzip(r *http.Request) bool {
	for _, part := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
		enc, params, _ := strings.Cut(strings.TrimSpace(part), ";")
		if strings.EqualFold(strings.TrimSpace(enc), "gzip") {
			return strings.ReplaceAll(params, " ", "") != "q=0"
		}
	}
	return false
}

// compress gzips HTML and CSV responses for clients that accept it. Responses that
// already carry a Content-Encoding (the precompressed assets) pass through.
func compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !acceptsGzip(r) {
			next.ServeHTTP(w, r)
			return
		}
		gw := &gzipWriter{ResponseWriter: w}
		defer gw.close()
		next.ServeHTTP(gw, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	zw          *gzip.Writer
	wroteHeader bool
}

func (g *gzipWriter) WriteHeader(code int) {
	if g.wroteHeader {
		return
	}
	g.wroteHeader = true
	h := g.Header()
	ct := h.Get("Content-Type")
	if code != http.StatusNoContent && code != http.StatusNotModified && h.Get("Content-Encoding") == "" &&
		(strings.HasPrefix(ct, "text/html") || strings.HasPrefix(ct, "text/csv")) {
		h.Set("Content-Encoding", "gzip")
		h.Add("Vary", "Accept-Encoding")
		h.Del("Content-Length")
		g.zw = gzip.NewWriter(g.ResponseWriter)
	}
	g.ResponseWriter.WriteHeader(code)
}

func (g *gzipWriter) Write(b []byte) (int, error) {
	if !g.wroteHeader {
		if g.Header().Get("Content-Type") == "" {
			g.Header().Set("Content-Type", http.DetectContentType(b))
		}
		g.WriteHeader(http.StatusOK)
	}
	if g.zw != nil {
		return g.zw.Write(b)
	}
	return g.ResponseWriter.Write(b)
}

func (g *gzipWriter) close() {
	if g.zw != nil {
		_ = g.zw.Close()
	}
}

// Unwrap lets http.ResponseController reach the underlying writer.
func (g *gzipWriter) Unwrap() http.ResponseWriter { return g.ResponseWriter }

// secure sets the security headers an on-prem deployment should send. Scripts and
// stylesheets are same-origin files. Components write a few inline style="…"
// attributes (progress width, aspect ratio, visually-hidden inputs), so
// style-src-attr needs 'unsafe-inline'; style elements and scripts stay strict.
func secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; style-src-attr 'unsafe-inline'; img-src 'self' data:; font-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}

// newMux registers the showcase's routes on top of the store.
func newMux(s *store) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", assetHandler()))

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/apps", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /apps", func(w http.ResponseWriter, r *http.Request) {
		vm := pages.ShowcaseVM{
			KPIs: []patterns.StatProps{
				{Label: "Deploys this week", Value: patterns.Int(128), Delta: patterns.SignedPct(18, 0), Good: true, Up: true},
				{Label: "Change failure rate", Value: patterns.Pct(2.3, 1), Delta: patterns.Signed(-0.8, 1) + " pp", Good: true, Up: false},
				{Label: "Mean time to restore", Value: patterns.Dec(41, 0) + " min", Delta: patterns.Signed(6, 0) + " min", Good: false, Up: true},
			},
			Apps:  s.search(r.URL.Query().Get("q")),
			Query: r.URL.Query().Get("q"),
		}
		render(w, r, http.StatusOK, pages.Showcase(vm))
	})

	// The Export button's download: every application, as CSV.
	mux.HandleFunc("GET /apps.csv", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", `attachment; filename="applications.csv"`)
		cw := csv.NewWriter(w)
		_ = cw.Write([]string{"application", "environment", "version", "health"})
		for _, a := range s.search("") {
			_ = cw.Write([]string{a.Name, a.Env, a.Version, a.Health})
		}
		cw.Flush()
		if err := cw.Error(); err != nil {
			log.Printf("apps.csv: %v", err)
		}
	})

	mux.HandleFunc("GET /apps/rows", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, http.StatusOK, pages.AppRows(s.search(r.URL.Query().Get("q"))))
	})

	mux.HandleFunc("POST /apps", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		e := pages.NewAppErrors{Values: map[string]string{"name": r.FormValue("name"), "replicas": r.FormValue("replicas")}}
		if n := r.FormValue("name"); n == "" || strings.ToLower(n) != n || strings.ContainsAny(n, " _.") {
			e.Name = "Use lowercase letters, digits and dashes."
		}
		if n, err := strconv.Atoi(r.FormValue("replicas")); err != nil || n < 2 {
			e.Replicas = "Must be at least 2 in production."
		}
		if e.Name != "" || e.Replicas != "" {
			render(w, r, http.StatusUnprocessableEntity, pages.NewAppForm(e))
			return
		}
		render(w, r, http.StatusOK, templ.Join(
			pages.NewAppForm(pages.NewAppErrors{}),
			toastPartial("Application created", r.FormValue("name")+" will sync within a minute.", toast.TypeSuccess),
		))
	})

	mux.HandleFunc("POST /apps/{name}/{action}", func(w http.ResponseWriter, r *http.Request) {
		name, action := r.PathValue("name"), r.PathValue("action")
		title := map[string]string{"restart": "Restarting pods", "sync": "Sync started"}[action]
		if title == "" {
			http.NotFound(w, r)
			return
		}
		render(w, r, http.StatusOK, toastPartial(title, name, toast.TypeInfo))
	})

	mux.HandleFunc("DELETE /apps/{name}", func(w http.ResponseWriter, r *http.Request) {
		s.remove(r.PathValue("name"))
		render(w, r, http.StatusOK, templ.Join(
			pages.AppRows(s.search(r.FormValue("q"))),
			toastPartial("Application deleted", r.PathValue("name"), toast.TypeSuccess),
		))
	})

	mux.HandleFunc("POST /scale", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, http.StatusOK, toastPartial("Pull request opened", "Scale to "+r.FormValue("replicas")+" replicas", toast.TypeSuccess))
	})

	mux.HandleFunc("GET /rollout", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		s.tick++
		ready := min(2+s.tick/3, 3)
		s.mu.Unlock()
		render(w, r, http.StatusOK, pages.RolloutCard(ready, 3, "every 2s"))
	})

	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, http.StatusOK, pages.Metrics(s.metrics(r)))
	})

	// The filter form's target: the body, plus the trigger label out of band.
	mux.HandleFunc("GET /metrics/body", func(w http.ResponseWriter, r *http.Request) {
		vm := s.metrics(r)
		render(w, r, http.StatusOK, templ.Join(
			pages.MetricsBody(vm),
			partial("#range-label", "outerHTML", pages.RangeLabel(vm.From, vm.To)),
		))
	})

	mux.HandleFunc("GET /metrics/release", func(w http.ResponseWriter, r *http.Request) {
		vm := s.metrics(r)
		app, version := r.URL.Query().Get("app"), r.URL.Query().Get("version")
		for _, rel := range vm.Releases {
			if rel.App == app && rel.Version == version {
				render(w, r, http.StatusOK, pages.ReleaseDetails(rel))
				return
			}
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("GET /settings", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, http.StatusOK, pages.Settings())
	})

	mux.HandleFunc("POST /settings/approve", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("code") != pages.ShowcaseApprovalCode {
			render(w, r, http.StatusUnprocessableEntity, pages.ApprovalForm("That code is not valid. Enter the current code from your authenticator.", false))
			return
		}
		render(w, r, http.StatusOK, templ.Join(
			pages.ApprovalForm("", true),
			toastPartial("Deployment freeze lifted", "Until 18:00 today", toast.TypeSuccess),
		))
	})

	// Everything unmatched: the Arion 404 page for page loads and boosted links.
	// htmx fragment requests keep the plain 404, so a broken hx-get can't swap a
	// whole page into a table body.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request-Type") == "partial" {
			http.NotFound(w, r)
			return
		}
		render(w, r, http.StatusNotFound, pages.NotFound(r.URL.Path))
	})

	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           secure(compress(newMux(newStore()))),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	// SIGTERM (OpenShift stopping the pod) or Ctrl-C: stop accepting, let requests
	// in flight finish, then exit.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	errc := make(chan error, 1)
	go func() {
		log.Printf("showcase on http://localhost:%s", port)
		errc <- srv.ListenAndServe()
	}()
	select {
	case err := <-errc:
		log.Fatal(err)
	case <-ctx.Done():
	}
	stop()
	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdown); err != nil {
		log.Printf("shutdown: %v", err)
	}
	log.Print("showcase stopped")
}
