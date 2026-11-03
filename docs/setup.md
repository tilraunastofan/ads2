# Setup

How to start a Go project on this system, and how apps take its updates. The design system is the Go module `github.com/tilraunastofan/ads2`; its `assets/css/arion-theme.css` is generated from this system's `tokens.json`, and the showcase in `cmd/showcase` is the reference app.

## Stack

| Layer | Choice | Notes |
| --- | --- | --- |
| Server | Go 1.26+, `net/http` `ServeMux` | Handlers render templ components; no JSON for UI routes. |
| Templates | templ (`github.com/a-h/templ`) | Pin as a Go tool: `go get -tool github.com/a-h/templ/cmd/templ`. Commit generated `*_templ.go`; CI runs `templ generate` and fails if `git status --porcelain -- '*_templ.go'` lists anything (a plain `git diff` misses new, untracked files). |
| Design system | `github.com/tilraunastofan/ads2` (this repo) | Imported as a Go module: components, patterns and assets (see **How apps use ads2**). |
| Components | shadcn-templ v2.0.0, style **nova** | Installed and themed in ads2 (`shadcn-templ add <name>` here, never in an app). |
| Interactivity | htmx 4 | Vendored in ads2's `assets/js/` and served by `assets.Handler` (on-prem: no runtime CDN). See **htmx patterns**. |
| CSS | Tailwind CSS v4.1+ standalone CLI | The app's `globals.css` → `output.css`, with ads2's classes through `go tool ads2css`. |
| Fonts | Hanken Grotesk, Source Serif 4, JetBrains Mono | Self-hosted variable woff2 in ads2's `assets/fonts/` (SIL OFL), loaded by `fonts.css`. No CDN. |

## How apps use ads2

An app **imports** the design system as a Go module, `github.com/tilraunastofan/ads2`. It never copies it, so `go get -u` brings every fix, token change and new pattern. One module, three parts that always move together:

| Package | What the app gets |
| --- | --- |
| `github.com/tilraunastofan/ads2/components/<name>` | The shadcn-templ components, themed and checked here. |
| `github.com/tilraunastofan/ads2/patterns` | `Shell`, `Base`, `PageHeader`, `Stat`, `DataCard`, number and date formatting … (`docs/patterns.md`). |
| `github.com/tilraunastofan/ads2/assets` | The theme CSS, fonts, htmx, the helper scripts and the component script bundle, embedded, with `assets.Handler` to serve them. |

The components come from the module too, not from `shadcn-templ add` in the app: `patterns` is built on these exact packages, and a second copy of a component would load its script twice. Need a component ads2 doesn't have? Add it here (`shadcn-templ add <name>` in ads2, with a doc in `docs/components/`), so every app gets it. Customise with `Class` and `Attributes`, never a fork.

Already have an app that copied `components/` and `patterns/`? Follow `docs/migrate-to-module.md`.

## New project

```sh
mkdir my-app && cd my-app
go mod init example.com/my-app
go get github.com/tilraunastofan/ads2@latest
go get -tool github.com/a-h/templ/cmd/templ@v0.3.1100     # the templ version ads2 is generated with
go get -tool github.com/tilraunastofan/ads2/cmd/ads2css   # connects Tailwind to the module
```

Then four files (the end-to-end smoke test in `testdata/consumer/` is exactly this app, and CI builds it):

1. **`assets/css/globals.css`**, the Tailwind entry:

   ```css
   @import "tailwindcss";
   @import "./ads2.gen.css";
   ```

   `go tool ads2css` writes `ads2.gen.css`: an `@import` of the module's `assets/css/ads2.css` (theme, shadcn styles) and `@source` lines for its components and patterns, so Tailwind generates the classes they use. The path is the module cache of the machine it ran on: **gitignore `ads2.gen.css` and `output.css`** and generate them in every build.

2. **The server** serves the design system's assets with the app's own on top, under `/assets/` (where `patterns.Base` links):

   ```go
   //go:embed static
   var static embed.FS // static/css/output.css, the app's Tailwind build, plus its own images

   app, _ := fs.Sub(static, "static")
   mux.Handle("GET /assets/", http.StripPrefix("/assets/", assets.Handler(assets.Overlay(app))))
   http.ListenAndServe(addr, patterns.SidebarState(mux))
   ```

   `assets.Overlay(app)` serves the app's file when it has one and the module's otherwise; `assets.Handler` adds content-hash ETags, gzip and a year-long cache for fonts and the hashed bundle. Keep app files out of the module's folders' names (`css/arion-theme.css`, `js/theme.js` …) unless you mean to override them.

3. **Pages** use `patterns.Shell` (or `Base`) and the module's components:

   ```templ
   import (
       "github.com/tilraunastofan/ads2/components/button"
       "github.com/tilraunastofan/ads2/patterns"
   )
   ```

4. **The build** (Taskfile or Makefile), in this order:

   ```sh
   go tool templ generate
   go tool ads2css -o static/css/ads2.gen.css        # next to globals.css
   tailwindcss -i static/css/globals.css -o static/css/output.css --minify
   go build ./...
   ```

Notes:

- Updating: `go get -u github.com/tilraunastofan/ads2 && go mod tidy`, then the build above (ads2css again: the module path changes with the version). Read `CHANGELOG.md` in ads2 for what changed and what an app has to do.
- templ: generate the app's templates with the templ version in ads2's `go.mod` (`go get -tool github.com/a-h/templ/cmd/templ@<that version>`); a mismatch prints a version warning on every render.
- Run `go tool templ generate && go mod tidy` after every `.templ` edit. Commit the generated `*_templ.go`.
- Working on ads2 and an app at once: `go mod edit -replace github.com/tilraunastofan/ads2=../ads2` in the app (ads2css follows the replace). Drop the replace before committing.
- Never edit the module's files in the module cache, and don't override `arion-theme.css`: theme changes are made in ads2's `tokens/tokens.json` (`task theme`), so every app gets them.
- Keep `body` class `style-nova` (`patterns.Base` sets it) and menus solid: do not enable the translucent menu colour.

## Theme variables

`arion-theme.css` defines exactly shadcn's variables, so every shadcn-templ component is themed without changes, plus these extra `--color-*` utilities: `ink`, `ink-foreground`, `ink-hover`, `slate`, `gold`, `plum`, `sand`, `positive`, `positive-foreground`, `warning`, `warning-foreground`, and `font-display` for the serif titles. It also defines the type scale (`text-display-sm`, `text-body-sm`, `text-code` …), maps `shadow-sm/md/lg` to the elevation tokens, and turns off motion for `prefers-reduced-motion`. Radii are Arion's 4 / 6 / 8 / 10px (`rounded-sm` … `rounded-xl`).

## Layout skeleton

Use `patterns.Shell` (sidebar, theme toggle, 1200px column) or `patterns.Base` (bare document) instead of writing the document by hand. With `Shell`, wrap the server's handler in `patterns.SidebarState(mux)` so a collapsed sidebar stays collapsed across pages. They load, in this order: `fonts.css`, `output.css`, `theme.js` (blocking, so dark mode has no flash), vendored htmx 4 (`defer`), the helpers `toast-guard.js`, `filter-form.js`, `request-errors.js`, `dialog-form.js`, `datepicker.js` and `focus-after-swap.js` (`defer`, each once: a `<script>` inside a page runs again on every boosted visit), `components.Scripts()`; and they mount `<div id="toasts" hidden>` (the target of `<hx-partial hx-target="#toasts">` responses) plus `toast.Toaster()`. The body is boosted with `hx-boost:inherited="true"`; `lang` defaults to `en` (`ShellProps.Lang` for Icelandic pages).

```templ
templ App(title, currentPath string) {
    @patterns.Shell(patterns.ShellProps{
        Title: title + " – Platform", Product: "Platform", Tagline: "On-prem OpenShift",
        GroupLabel: "Workloads", Nav: nav(), CurrentPath: currentPath,
    }) {
        { children... }
    }
}
```

## Verifying UI

- Dark mode is the class `dark` on `<html>`, which swaps the token variables and enables `dark:`. `theme.js` sets it from the OS setting on first visit and from the toggle's saved choice afterwards. The CSS itself ignores `prefers-color-scheme`.
- Check toasts, dialogs, menus and themes in a real browser; HTTP checks only cover the server side.
- Browser extensions such as Dark Reader repaint pages. Judge the theme by its variables (`getComputedStyle(document.documentElement).getPropertyValue('--background')`), not by screenshot colour.
- `task dev` needs ports 8090 and 7331 free (`lsof -nP -iTCP:8090 -sTCP:LISTEN`).

## On-prem checklist

The showcase (`cmd/showcase`) implements all of it; take the parts that are app code (CSP, timeouts, shutdown) from there.

- htmx, the fonts and the scripts are embedded in the module and in the app's binary through `assets.Handler`: no CDN, nothing to download at runtime.
- Caching is `assets.Handler`'s: fonts and hashed `js/shadcn-templ-*.js` bundles `immutable` for a year, everything else `no-cache` with an ETag.
- Container: build with `task build`, run as non-root on OpenShift (random UID): write nothing to the image filesystem at runtime.
- CSP: scripts and styles are files, so `script-src 'self'; style-src 'self'` works. Components write a few inline `style="…"` attributes (progress width, aspect ratio, visually hidden inputs), so also send `style-src-attr 'unsafe-inline'`; this does not allow `<style>` elements. If you use `hx-on:*`, add `htmx.config.inlineScriptNonce`. htmx trigger filters (`hx-trigger="input[target.value.length==6]"`) are compiled with `new Function`, which `script-src 'self'` blocks without any visible console error; use a plain event name and filter on the server. `input-otp` also injects a `<style>` at runtime, which `style-src 'self'` blocks (its autofill and iOS rules are dropped; typing still works). The showcase's `secure` middleware sends a working policy.
