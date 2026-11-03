package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Handler serves fsys over HTTP, for mounting under /assets/, the path patterns.Base
// links to:
//
//	mux.Handle("GET /assets/", http.StripPrefix("/assets/", assets.Handler(assets.Overlay(app))))
//
// Fonts and the hashed component bundle (js/shadcn-templ-*.js) never change under
// their name, so they are cached for a year. Everything else revalidates: embedded
// files carry no modification time, so each gets a content-hash ETag and a refresh
// costs a 304. Text types (css, js, svg, json, txt) are gzipped once, on first
// request, for clients that accept gzip, with their own ETag.
//
// Files are read once and kept in memory, so fsys must not change while serving (an
// embed.FS never does). Serve a directory being edited with http.FileServer instead.
func Handler(fsys fs.FS) http.Handler {
	var cache sync.Map // name -> *cachedAsset
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/")
		w.Header().Set("Cache-Control", CacheControl(name))
		a, ok := cache.Load(name)
		if !ok {
			body, err := fs.ReadFile(fsys, name)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			a, _ = cache.LoadOrStore(name, newCachedAsset(name, body))
		}
		ca := a.(*cachedAsset)
		body, etag := ca.body, ca.etag
		if ca.gz != nil {
			w.Header().Add("Vary", "Accept-Encoding")
			if AcceptsGzip(r) {
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

// CacheControl is the Cache-Control value Handler sends for an asset name (relative
// to /assets/): a year for fonts and the hashed component bundle, revalidation for
// everything else.
func CacheControl(name string) string {
	if strings.HasPrefix(name, "fonts/") || strings.HasPrefix(name, "js/shadcn-templ-") {
		return "public, max-age=31536000, immutable"
	}
	return "no-cache"
}

// Overlay is an app's own assets on top of the design system's: a file in app wins,
// any other comes from Assets. An app serves its Tailwind build (css/output.css) and
// its own images or scripts this way, under the same /assets/ and with the same
// caching, while the theme, fonts, htmx, helpers and component bundle come from the
// ads2 version in its go.mod.
//
//	//go:embed static
//	var static embed.FS
//	app, _ := fs.Sub(static, "static") // static/css/output.css → css/output.css
//	assets.Overlay(app)
func Overlay(app fs.FS) fs.FS { return overlay{app, Assets} }

type overlay struct{ top, base fs.FS }

func (o overlay) Open(name string) (fs.File, error) {
	f, err := o.top.Open(name)
	if err == nil {
		return f, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	return o.base.Open(name)
}

// AcceptsGzip reports whether the request accepts a gzip response: "gzip" in
// Accept-Encoding with a q-value above zero ("gzip;q=0" and "gzip;q=0.0" refuse it).
func AcceptsGzip(r *http.Request) bool {
	for part := range strings.SplitSeq(r.Header.Get("Accept-Encoding"), ",") {
		enc, params, _ := strings.Cut(strings.TrimSpace(part), ";")
		if !strings.EqualFold(strings.TrimSpace(enc), "gzip") {
			continue
		}
		k, v, ok := strings.Cut(strings.ReplaceAll(params, " ", ""), "=")
		if !ok || !strings.EqualFold(k, "q") {
			return true
		}
		q, err := strconv.ParseFloat(v, 64)
		return err == nil && q > 0
	}
	return false
}

// cachedAsset is a file's bytes and ETag and, for text types, its gzipped bytes,
// computed once per process.
type cachedAsset struct {
	body, gz     []byte
	etag, gzETag string
}

func newCachedAsset(name string, body []byte) *cachedAsset {
	sum := sha256.Sum256(body)
	h := hex.EncodeToString(sum[:8])
	a := &cachedAsset{body: body, etag: `"` + h + `"`}
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
