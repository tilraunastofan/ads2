// Package main is the smallest app on the ads2 module, the one docs/setup.md
// describes. scripts/consumer-smoke.sh builds and serves it against this checkout,
// so CI proves that an app can import ads2. It has no go.mod here (testdata is not
// built with the repo); the script writes one.
package main

import (
	"cmp"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"example.com/consumer/ui"
	"github.com/tilraunastofan/ads2/assets"
	"github.com/tilraunastofan/ads2/patterns"
)

//go:embed static
var static embed.FS

func main() {
	app, _ := fs.Sub(static, "static")
	mux := http.NewServeMux()
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", assets.Handler(assets.Overlay(app))))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) { _ = ui.Home().Render(r.Context(), w) })
	addr := ":" + cmp.Or(os.Getenv("PORT"), "8099")
	log.Fatal(http.ListenAndServe(addr, patterns.SidebarState(mux)))
}
