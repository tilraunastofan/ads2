// Command showcase serves the Arion Design System 2 reference app: shadcn-templ
// components themed with assets/css/arion-theme.css, driven by htmx 4.
//
//	task dev          # templ + tailwind + script watchers, http://localhost:7331
//	PORT=8090 go run ./cmd/showcase
package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/a-h/templ"

	"github.com/tilraunastofan/ads2/components/toast"
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

func main() {
	s := newStore()
	mux := http.NewServeMux()

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/apps", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /apps", func(w http.ResponseWriter, r *http.Request) {
		vm := pages.ShowcaseVM{
			KPIs: []pages.KPI{
				{Label: "Deploys this week", Value: "128", Delta: "+18%", Good: true, Up: true},
				{Label: "Change failure rate", Value: "2,3%", Delta: "−0,8 pp", Good: true, Up: false},
				{Label: "Mean time to restore", Value: "41 mín", Delta: "+6 mín", Good: false, Up: true},
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
		render(w, r, http.StatusOK, pages.RolloutCard(ready, 3))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	log.Printf("showcase on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
