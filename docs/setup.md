# Setup

How to start a Go project on this system. The reference starter lives in the `ads2` repository; its `assets/css/arion-theme.css` is generated from this system's `tokens.json`.

## Stack

| Layer | Choice | Notes |
| --- | --- | --- |
| Server | Go 1.26+, `net/http` `ServeMux` | Handlers render templ components; no JSON for UI routes. |
| Templates | templ (`github.com/a-h/templ`) | Pin as a Go tool: `go get -tool github.com/a-h/templ/cmd/templ`. Commit generated `*_templ.go`; CI runs `templ generate` + `git diff --exit-code`. |
| Components | shadcn-templ v2 (`github.com/axadrn/shadcn-templ/v2`) | Style **nova** (`base-nova` in `components.json`). Add components with `shadcn-templ add <name>`. |
| Interactivity | htmx 4 | Vendored into `assets/js/` (on-prem: no runtime CDN). See **htmx patterns**. |
| CSS | Tailwind CSS v4.1+ standalone CLI | `assets/css/globals.css` → `assets/css/output.css`. |
| Fonts | Hanken Grotesk, Source Serif 4, JetBrains Mono | Self-hosted variable woff2 in `assets/fonts/` (SIL OFL), loaded by `assets/css/fonts.css`. No CDN. |

## New project

```sh
go install github.com/axadrn/shadcn-templ/v2/cmd/shadcn-templ@latest
shadcn-templ init my-app -t templ --preset nova
cd my-app
shadcn-templ add button card badge table field input input-group textarea label checkbox radio-group switch native-select select slider toggle toggle-group tabs breadcrumb pagination accordion sidebar avatar item kbd separator alert toast progress spinner skeleton empty dialog alert-dialog dropdown-menu tooltip icon
go mod tidy && go tool templ generate && go build ./... && go vet ./...
```

- Registry names are **kebab-case** (`alert-dialog`, `dropdown-menu`, `input-group`, `native-select`); the Go package names drop the hyphen (`alertdialog`, `dropdownmenu`…). The no-hyphen names return 404.
- `input-group` needs `textarea` installed too.
- In a fresh scaffold run `go mod tidy` first: `go tool templ` is declared in `go.mod` but only resolves after tidy. After that, run `go tool templ generate && go mod tidy` after every `.templ` edit or `shadcn-templ add`.
- Never edit the vendored `components/` or the generated `arion-theme.css`: change `tokens.json` and regenerate (`task theme`).

Then apply the Arion theme:

1. Copy `arion-theme.css`, `fonts.css`, the `assets/fonts/` folder, `assets/js/theme.js`, `assets/img/arion-mark.svg` and the `patterns/` package from `ads2` into the project (rewrite the `patterns` import paths to your module).
2. In `globals.css` delete the `:root { … }`, `.dark { … }`, `@theme inline { … }` and `@custom-variant dark` blocks that `init` wrote, and add `@import "./arion-theme.css";` after the other `@import` lines.
3. Keep `body` class `style-nova` (the init scaffold sets it). Keep menus solid: do not enable the translucent menu colour.
4. Link `fonts.css` before `output.css` (`patterns.Base` does this).

## Theme variables

`arion-theme.css` defines exactly shadcn's variables, so every shadcn-templ component is themed without changes, plus these extra `--color-*` utilities: `ink`, `ink-foreground`, `ink-hover`, `slate`, `gold`, `plum`, `sand`, `positive`, `positive-foreground`, `warning`, `warning-foreground`, and `font-display` for the serif titles. It also defines the type scale (`text-display-sm`, `text-body-sm`, `text-code` …), maps `shadow-sm/md/lg` to the elevation tokens, and turns off motion for `prefers-reduced-motion`. Radii are Arion's 4 / 6 / 8 / 10px (`rounded-sm` … `rounded-xl`).

## Layout skeleton

Use `patterns.Shell` (sidebar, theme toggle, 1200px column) or `patterns.Base` (bare document) instead of writing the document by hand. They load, in this order: `fonts.css`, `output.css`, `theme.js` (blocking, so dark mode has no flash), vendored htmx 4 (`defer`), `components.Scripts()`; and they mount `<div id="toasts" hidden>` (the target of `<hx-partial hx-target="#toasts">` responses) plus `toast.Toaster()`. The body is boosted with `hx-boost:inherited="true"`; `lang` defaults to `en` (`ShellProps.Lang` for Icelandic pages).

```templ
templ App(title, currentPath string) {
    @patterns.Shell(patterns.ShellProps{
        Title: title + " · Platform", Product: "Platform", Tagline: "on-prem · OpenShift",
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

All of this is implemented in `cmd/showcase`; copy it.

- Vendor htmx (`assets/js/htmx.min.js`, 4.0.x) and the fonts (`assets/fonts/`). Embed `/assets/` with `go:embed` (`assets/assets.go`) and serve it from the binary; read from disk when `SHADCN_TEMPL_DEV=true`.
- Hashed `assets/js/shadcn-templ-*.js` bundles: `Cache-Control: immutable`. Everything else `no-cache`.
- Container: build with `task build`, run as non-root on OpenShift (random UID): write nothing to the image filesystem at runtime.
- CSP: scripts and styles are files, so `script-src 'self'; style-src 'self'` works. Components write a few inline `style="…"` attributes (progress width, aspect ratio, visually hidden inputs), so also send `style-src-attr 'unsafe-inline'`; this does not allow `<style>` elements. If you use `hx-on:*`, add `htmx.config.inlineScriptNonce`. The showcase's `secure` middleware sends a working policy.
