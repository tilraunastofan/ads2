# ads2 — Arion Design System 2

Arion banki's palette on **shadcn-templ v2** (nova style) for Go + templ + htmx 4 + Tailwind v4.
Module: `github.com/tilraunastofan/ads2`.

## Layout

| Path | What |
| --- | --- |
| `tokens/tokens.json` | Source of truth for colours, type, spacing, radii (exported from the design system artifact). |
| `cmd/gentheme` | Generates `assets/css/arion-theme.css` (colours, type scale, shadows, radii, reduced motion) from the tokens and fails if a contrast requirement in `contrast.go` is missed. Never edit that CSS by hand. |
| `assets/css/arion-theme.css` | shadcn theme variables (`:root`, `.dark`, `@theme inline`) + Arion extras. |
| `assets/css/fonts.css`, `assets/fonts/` | Self-hosted variable woff2 (OFL). Suisse is optional and opt-in (see Typography in `design-system.md`). |
| `assets/img/arion-mark.svg` | Arion star mark (vector trace of the Arion logo; use as is). |
| `assets/js/theme.js` | Dark-mode class on `<html>`: OS setting first, then the toggle's saved choice. |
| `assets/js/datepicker.js` | Copies a calendar selection into hidden inputs of a `form[data-date-picker]` and fires `change` (see `docs/components/DatePicker.md`). |
| `assets/assets.go` | `go:embed` of `css`, `fonts`, `img`, `js` for the binary. |
| `patterns/` | Reusable compositions: `Shell`, `Base`, `PageHeader`, `Stat`, `StatusBadge`, `DataCard`, `ThemeToggle`, number and date formatting (`format.go`, `dates.go`). Copied into apps like `components/`. |
| `ui/layouts/app.templ` | The showcase's navigation on top of `patterns.Shell`. |
| `ui/pages/showcase.templ` | Reference page using the patterns and ~20 components with htmx patterns. |
| `ui/pages/metrics.templ` | `/metrics`: DatePicker (popover + range calendar), Combobox filter, Chart, DataCard with a Sheet per release. |
| `ui/pages/settings.templ` | `/settings`: step-up approval with InputOTP (showcase code `123456`). |
| `.github/workflows/ci.yml` | `task check`, then `task build` with the Tailwind standalone CLI and a smoke test of `/apps`, `/metrics`, `/settings`. |
| `cmd/showcase` | `net/http` server for the showcase: embedded assets, cache headers, CSP. |
| `docs/` | Brand book (`design-system.md`), `setup.md`, `htmx-patterns.md`, `patterns.md`, one file per component. **Read these before writing UI.** |

## Setup & commands

Prerequisites: Go ≥ 1.26 (verified on 1.27), [Task](https://taskfile.dev), Tailwind CSS v4.1+ standalone CLI (`brew install tailwindcss`).

Already done (committed): `shadcn-templ init --preset nova --css assets/css/globals.css`, `globals.css` imports `arion-theme.css` instead of defining its own theme blocks, htmx 4.0.0 vendored at `assets/js/htmx.min.js`, all components named in `docs/components/` installed, generated `*_templ.go` and `components/scripts_bundle.go`.

- Add components with the registry's **kebab-case** names: `shadcn-templ add alert-dialog dropdown-menu input-group native-select …` (`alertdialog`, `dropdownmenu`, `inputgroup`, `nativeselect` 404). `input-group` needs `textarea` added explicitly.
- `task dev` — tailwind + script bundler + templ watcher; open http://localhost:7331/apps (app itself listens on :8090, override with `PORT`, but the templ proxy in `Taskfile.yml` targets :8090).
- `task build` — production build to `bin/showcase`.
- `task check` — CI: theme current and contrast met, templ output and script bundle committed, `go vet`, `go build`, `go test`.
- `task theme` — regenerate `assets/css/arion-theme.css` after changing `tokens/tokens.json`; it also checks every contrast claim and fails with the offending pair.
- `go test ./...` — token contrast, theme drift (`arion-theme.css` must match `tokens.json`), number formatting.
- After adding components or editing `.templ`: `go tool templ generate && go mod tidy`.

Verified in Chrome against `task dev`: search filters rows, row menu raises a toast, delete asks in an alert dialog and removes the row, the New application form shows 422 field errors, Rollout reaches 3/3 and stops polling, dark mode via `class="dark"` on `<html>` swaps the tokens. `/metrics`: the range picker updates the label, KPIs, chart and table, the app filter narrows them, Details opens a sheet whose body loads lazily. `/settings`: `111111` shows the 422 error, `123456` shows "Freeze lifted" and a toast. Note: a dark-mode browser extension (e.g. Dark Reader) darkens the page even without `.dark`; disable it when judging themes.

## Notes for agents

**Tooling**
- `go tool templ generate && go mod tidy`, then `go build ./... && go vet ./...`, after any `.templ` edit or `shadcn-templ add`.
- Don't edit `components/` or `assets/css/arion-theme.css`; change `tokens/tokens.json` and run `task theme`.
- `task dev` needs :8090 and :7331 free. If it fails to bind, check `lsof -nP -iTCP:8090 -sTCP:LISTEN` before assuming a code problem, and don't kill processes you didn't start.
- Charts: set colours on the series (`Fill: "var(--chart-1)"`), never `chart.Config.Color`; a config colour makes the container write an inline `<style>`, which the CSP (`style-src 'self'`) blocks.
- Calendar and DatePicker: weeks start Monday (`calendar.Monday`); `Locale` `en-GB` on English pages, `is` on Icelandic ones. Dates in text come from `patterns.Date` / `DateRange`.
- htmx polling: `hx-trigger="every 2s"` starts at page load even inside a hidden tab, and `intersect once, every 2s` does too (`every` is its own trigger). To start polling when a tab opens, render `hx-trigger="intersect once"` first and have the response carry `every 2s` (see `RolloutCard`).
- Run `templ generate` only when `task dev` is stopped: the watcher regenerates on save, and a manual run deletes its temp files, after which the :7331 proxy serves a blank page (:8090 still works; restart `task dev`).
- The Tailwind watcher only rebuilds on a change. After pulling or switching branches with new templates, run `tailwindcss -i ./assets/css/globals.css -o ./assets/css/output.css` once; otherwise new classes (`h-[240px]`) are missing and e.g. charts collapse to zero height.
- CSP blocks `eval`: htmx trigger filters (`input-otp-change[detail.value.length==6]`) never fire and log nothing in the console. Use plain event triggers. Listen for `securitypolicyviolation` to see what is blocked.
- The showcase store is in memory; a restart brings deleted rows back.

**Verifying UI**
- Check toasts, dialogs and dark mode in Chrome; HTTP only covers server responses. Check `list_connected_browsers` if the extension reports not connected.
- Dark Reader (if installed) darkens pages even without `.dark`. Judge themes by reading CSS variables, not screenshot colours.
- Dark tokens and the `dark:` variant both follow `.dark` on `<html>` (not `prefers-color-scheme`, not `data-theme`); `theme.js` sets the class from the OS setting and the saved toggle.

## Rules for UI work in this repo

- Use shadcn-templ components first; customise via `Class` and `Attributes` (htmx attributes go in `Attributes`). Don't fork components.
- Colours only through tokens (`bg-primary`, `text-muted-foreground`, `text-positive`, `bg-ink`…). No hex values or Tailwind palette colours (`blue-600`) in templates.
- One `primary` button per view. Destructive actions go through `alertdialog`, never `hx-confirm`.
- Numbers: Arion format in every UI language (`.` thousands, `,` decimal; use `patterns.Int/Dec/Signed/Pct`), `font-mono tabular-nums`, right-aligned; gains `text-positive` with `+`, losses `text-destructive` with `−`. Words and units follow the page language ("41 min" on English pages).
- State pills: `patterns.StatusBadge`, never ad-hoc badge variants. Page chrome: `patterns.Shell`, `PageHeader`, `Stat`, `DataCard` before hand-rolling.
- Type: use the scale utilities (`font-display text-display-sm` for the page title, `text-body-sm`, `text-code`…), never `text-[40px]`. Everything but titles is sans.
- htmx 4: inheritance is explicit (`hx-x:inherited`), 4xx/5xx swap by default (use `hx-status:422`), use `<hx-partial>` for multi-target responses (toasts go to `#toasts`).
- After changing `tokens/tokens.json`: `task theme`.
