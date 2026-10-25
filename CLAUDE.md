# ads2 — Arion Design System 2

Arion banki's palette on **shadcn-templ v2** (nova style) for Go + templ + htmx 4 + Tailwind v4.
Module: `github.com/tilraunastofan/ads2`.

## Layout

| Path | What |
| --- | --- |
| `tokens/tokens.json` | Source of truth for colours, type, spacing, radii (exported from the design system artifact). |
| `cmd/gentheme` | Generates `assets/css/arion-theme.css` from the tokens. Never edit that CSS by hand. |
| `assets/css/arion-theme.css` | shadcn theme variables (`:root`, `.dark`, `@theme inline`) + Arion extras. |
| `assets/css/fonts.css` | Font loading (Google in dev; swap for self-hosted woff2 on-prem). |
| `assets/img/arion-mark.svg` | Arion star mark (from `logos/logo.svg` / the design system assets). |
| `ui/layouts/app.templ` | `Base` document + `App` sidebar shell. |
| `ui/pages/showcase.templ` | Reference page using ~20 components with htmx patterns. |
| `cmd/showcase` | `net/http` server for the showcase. |
| `docs/` | Brand book (`design-system.md`), `setup.md`, `htmx-patterns.md`, one file per component. **Read these before writing UI.** |

## Setup & commands

Prerequisites: Go ≥ 1.26, [Task](https://taskfile.dev), Tailwind CSS v4.1+ standalone CLI (`brew install tailwindcss`).

Already done (committed): `shadcn-templ init --preset nova --css assets/css/globals.css`, `globals.css` imports `arion-theme.css` instead of defining its own theme blocks, htmx 4.0.0 vendored at `assets/js/htmx.min.js`, generated `*_templ.go` and `components/scripts_bundle.go`.

- Add components with the registry's **kebab-case** names: `shadcn-templ add alert-dialog dropdown-menu input-group native-select …` (`alertdialog`, `dropdownmenu`, `inputgroup`, `nativeselect` 404). `input-group` needs `textarea` added explicitly.
- `task dev` — tailwind + script bundler + templ watcher; open http://localhost:7331/apps (app itself listens on :8090, override with `PORT`, but the templ proxy in `Taskfile.yml` targets :8090).
- `task build` — production build to `bin/showcase`.
- `task check` — CI: templ output committed, `go vet`, `go build`.
- `task theme` — regenerate `assets/css/arion-theme.css` after changing `tokens/tokens.json`.
- After adding components or editing `.templ`: `go tool templ generate && go mod tidy`.

Verified in Chrome against `task dev`: search filters rows, row menu raises a toast, delete asks in an alert dialog and removes the row, the New application form shows 422 field errors, Rollout reaches 3/3 and stops polling, dark mode via `class="dark"` on `<html>` swaps the tokens. Note: a dark-mode browser extension (e.g. Dark Reader) darkens the page even without `.dark`; disable it when judging themes.

## Notes for agents

**Tooling**
- `go tool templ generate && go mod tidy`, then `go build ./... && go vet ./...`, after any `.templ` edit or `shadcn-templ add`.
- Don't edit `components/` or `assets/css/arion-theme.css`; change `tokens/tokens.json` and run `task theme`.
- `task dev` needs :8090 and :7331 free. If it fails to bind, check `lsof -nP -iTCP:8090 -sTCP:LISTEN` before assuming a code problem, and don't kill processes you didn't start.
- htmx polling: `hx-trigger="every 2s"` starts at page load even inside a hidden tab, and `intersect once, every 2s` does too (`every` is its own trigger). To start polling when a tab opens, render `hx-trigger="intersect once"` first and have the response carry `every 2s` (see `RolloutCard`).
- The showcase store is in memory; a restart brings deleted rows back.

**Verifying UI**
- Check toasts, dialogs and dark mode in Chrome; HTTP only covers server responses. Check `list_connected_browsers` if the extension reports not connected.
- Dark Reader (if installed) darkens pages even without `.dark`. Judge themes by reading CSS variables, not screenshot colours.
- Dark tokens apply only under `.dark` (not `prefers-color-scheme`); `[data-theme="dark"]` only enables the `dark:` variant and does not change the tokens.

## Rules for UI work in this repo

- Use shadcn-templ components first; customise via `Class` and `Attributes` (htmx attributes go in `Attributes`). Don't fork components.
- Colours only through tokens (`bg-primary`, `text-muted-foreground`, `text-positive`, `bg-ink`…). No hex values or Tailwind palette colours (`blue-600`) in templates.
- One `primary` button per view. Destructive actions go through `alertdialog`, never `hx-confirm`.
- Numbers: Icelandic format, `font-mono tabular-nums`, right-aligned; gains `text-positive` with `+`, losses `text-destructive` with `−`.
- Page titles: `font-display` (serif) 40px; everything else sans.
- htmx 4: inheritance is explicit (`hx-x:inherited`), 4xx/5xx swap by default (use `hx-status:422`), use `<hx-partial>` for multi-target responses (toasts go to `#toasts`).
- After changing `tokens/tokens.json`: `task theme`.
