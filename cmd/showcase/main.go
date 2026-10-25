// Command showcase serves the Arion Design System 2 reference app: shadcn-templ
// components themed with assets/css/arion-theme.css, driven by htmx 4.
//
//	task dev          # templ + tailwind + script watchers, http://localhost:7331
//	PORT=8090 go run ./cmd/showcase
package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
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
	var etags sync.Map // asset name -> ETag
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
		tag, ok := etags.Load(name)
		if !ok {
			sum := sha256.Sum256(body)
			tag, _ = etags.LoadOrStore(name, `"`+hex.EncodeToString(sum[:8])+`"`)
		}
		w.Header().Set("ETag", tag.(string))
		http.ServeContent(w, r, path.Base(name), time.Time{}, bytes.NewReader(body))
	})
}

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

func main() {
	s := newStore()
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	log.Printf("showcase on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, secure(mux)))
}
