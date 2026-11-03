# ads2 — Arion Design System 2

Arion banki's palette on **shadcn-templ v2.0.0** (nova style; the CLI is pinned as a Go tool in `go.mod`) for Go + templ + htmx 4 + Tailwind v4.
Module: `github.com/tilraunastofan/ads2`.

## Layout

| Path | What |
| --- | --- |
| `tokens/tokens.json` | Source of truth for colours, type, spacing, radii (exported from the design system artifact). |
| `cmd/gentheme` | Generates `assets/css/arion-theme.css` (colours, type scale, shadows, radii, reduced motion) from the tokens and fails if a contrast requirement in `contrast.go` is missed. Never edit that CSS by hand. |
| `assets/css/arion-theme.css` | shadcn theme variables (`:root`, `.dark`, `@theme inline`) + Arion extras. |
| `assets/css/fonts.css`, `assets/fonts/` | Self-hosted variable woff2 (OFL). Suisse is optional and opt-in (see Typography in `design-system.md`). |
| `assets/img/arion-mark.svg` | Arion star mark (vector trace of the Arion logo). In pages use `patterns.ArionMark` (inline, `brand-mark` token: blue in light, white in dark); `arion-mark-favicon.svg` is the favicon. |
| `assets/js/theme.js` | Dark-mode class on `<html>`: OS setting first, then the toggle's saved choice. |
| `assets/js/toast-guard.js` | Workaround for a freeze in the vendored toast: disconnects each toast's observers before a boosted body swap. Remove once shadcn-templ's toast `destroy()` does it ([shadcn-templ#656](https://github.com/axadrn/shadcn-templ/issues/656)). Relies on toast internals: after every shadcn-templ update, run the checks in `docs/known-issues/toast-freeze.md`. |
| `assets/js/filter-form.js` | `form[data-filter-form]` fires `filter-change` only when its named values change (see **Filter form** in `docs/htmx-patterns.md`). |
| `assets/js/request-errors.js` | Error toast for unreachable server and unhandled 5xx (cancels swapping raw error text); en/is wording. |
| `assets/js/dialog-form.js` | `form[data-close-on-success]` inside a dialog or sheet closes it after a 2xx (closing on the submit click cancels the submit). |
| `assets/js/datepicker.js` | Copies a calendar selection into hidden inputs of a `form[data-date-picker]` and fires `change` (see `docs/components/DatePicker.md`). |
| `assets/js/focus-after-swap.js` | When a swap removes the focused element: focus the new page's `<h1>` after a boosted navigation, or the first `aria-invalid="true"` field of a re-rendered form. All helpers load once in `patterns.Base`; never add a helper `<script>` to a page (it reruns on every boosted visit). |
| `assets/assets.go`, `assets/handler.go` | `go:embed` of `css`, `fonts`, `img`, `js`, and `assets.Handler` / `assets.Overlay`, which the showcase and every app serve them with. `assets/css/ads2.css` is the CSS entry apps import (theme, shadcn styles, base layer); this repo's `globals.css` imports it too. |
| `cmd/ads2css` | `go tool ads2css` in an app: writes `ads2.gen.css`, which imports the module's `ads2.css` and adds `@source` lines for its components and patterns. |
| `patterns/` | Reusable compositions: `Shell`, `Base`, `ReportShell`, `ReportCover`, `PageHeader`, `Stat`, `StatusBadge`, `DataCard`, `ThemeToggle`, number and date formatting (`format.go`, `dates.go`), Arion numbers for charts (`chart.templ`). Apps import it from the module, with `components/` and `assets/`. |
| `ui/layouts/app.templ` | The showcase's navigation on top of `patterns.Shell`. |
| `ui/pages/showcase.templ` | Reference page using the patterns and ~20 components with htmx patterns. |
| `ui/pages/metrics.templ` | `/metrics`: DatePicker (popover + range calendar), Combobox filter, Chart, DataCard with a Sheet per release. |
| `ui/pages/report.templ` | `/uppgjor`: Icelandic reference report on `patterns.ReportShell` + `ReportCover` (demo income statement, m.kr., Icelandic dates and chrome). |
| `ui/pages/settings.templ` | `/settings`: step-up approval with InputOTP (showcase code `123456`). |
| `.github/workflows/ci.yml` | `task check`, then `task build` with the Tailwind standalone CLI and a smoke test of `/apps`, `/metrics`, `/settings`, `/uppgjor`, the 404 status and gzip. |
| `.github/dependabot.yml` | Weekly grouped PRs for Go modules and GitHub Actions. Vendored htmx and the pinned Tailwind CLI are bumped by hand. |
| `cmd/showcase` | `net/http` server for the showcase: embedded assets (text types gzipped once at first request), cache headers, CSP, gzip for HTML/CSV, server timeouts, graceful shutdown on SIGTERM, panic recovery (toast for htmx requests, error page otherwise), Arion 404 page (`ui/pages/notfound.templ`; htmx fragment requests keep a plain 404). |
| `CHANGELOG.md` | Per release: what changed and what an app has to do. Apps read it when they `go get -u`. |
| `docs/migrate-to-module.md` | For agents: convert an app that copied `components/`, `patterns/` and the assets to importing the module. |
| `testdata/consumer`, `scripts/consumer-smoke.sh` | The smallest app on the module (the one `docs/setup.md` describes); CI builds it against each commit and checks every asset it links to. |
| `docs/known-issues/` | Bugs in vendored code we work around, with reproduction, cause and the checks to run after an upgrade (`toast-freeze.md`). |
| `docs/` | Brand book (`design-system.md`), `setup.md`, `htmx-patterns.md`, `patterns.md`, one file per component. **Read these before writing UI.** |

## Setup & commands

Prerequisites: Go ≥ 1.26 (verified on 1.27), [Task](https://taskfile.dev), Tailwind CSS v4.1+ standalone CLI (`brew install tailwindcss`).

Already done (committed): `shadcn-templ init --preset nova --css assets/css/globals.css`, `globals.css` imports `ads2.css` (which imports `arion-theme.css`) instead of defining its own theme blocks, htmx 4.0.1 vendored at `assets/js/htmx.min.js`, all components named in `docs/components/` installed, generated `*_templ.go` and `components/scripts_bundle.go`.

- Add components with the registry's **kebab-case** names: `shadcn-templ add alert-dialog dropdown-menu input-group native-select …` (`alertdialog`, `dropdownmenu`, `inputgroup`, `nativeselect` 404). `input-group` needs `textarea` added explicitly.
- `task dev` — tailwind + script bundler + templ watcher; open http://localhost:7331/apps (app itself listens on :8090, override with `PORT`, but the templ proxy in `Taskfile.yml` targets :8090).
- `task build` — production build to `bin/showcase`.
- `task check` — CI: theme current and contrast met, templ output and script bundle committed, `go vet`, `go build`, `go test -race`.
- `task theme` — regenerate `assets/css/arion-theme.css` after changing `tokens/tokens.json`; it also checks every contrast claim and fails with the offending pair.
- `go test ./...` — token contrast, theme drift (`arion-theme.css` must match `tokens.json`), number formatting, and `ui/designrules_test.go`, which fails on any `.templ` in `ui/` or `patterns/` that breaks the UI rules below (hex or palette colours, `text-[…]` and Tailwind text sizes, `hx-confirm`, raw `<button>`/`<table>`/`<input>`, an `<h1>` outside `PageHeader`, a Go bool on an `aria-*` attribute), and `TestOnePrimaryButtonPerView` (`cmd/showcase`), which renders each page and allows one primary button in the page plus any one tab panel, and one per dialog, sheet or menu.
- After adding components or editing `.templ`: `go tool templ generate && go mod tidy`.

Verified in Chrome against `task dev`: search filters rows; the row menu's Restart and Sync raise a toast, Scale… opens a dialog for that row that closes after the pull-request toast, Delete… asks in an alert dialog and removes the row; New application opens a sheet that keeps 422 field errors and on success closes, adds the row and toasts; Rollout shows "2 of 3", reaches "3 of 3" and stops polling; dark mode via `class="dark"` on `<html>` swaps the tokens. `/metrics`: the range picker updates the label, KPIs, chart and table, the app filter narrows them, Details opens a sheet whose body loads lazily, also right after changing the range and typing in the app filter, and shows the release time of its row. A toast showing during a boosted navigation or browser Back no longer freezes the tab. `/settings` (Deployment freeze): `111111` shows the 422 error, **Lift freeze** with `123456` shows "Freeze lifted" and a toast. Dark mode: the Arion mark turns white. Note: a dark-mode browser extension (e.g. Dark Reader) darkens the page even without `.dark`; disable it when judging themes.

## Notes for agents

**This repo is a Go module that apps import** (`patterns`, `components`, `assets`, `utils`, `cmd/ads2css`). An exported name, a prop, a CSS class a pattern relies on, an asset path or a token is API: apps get the change on their next `go get -u`.
- Breaking one is allowed before v1, but add it under the next version in `CHANGELOG.md` (**Breaking**, with the one thing an app does about it) in the same commit.
- Keep those packages free of `ui/` and `cmd/` imports (`TestModulePackagesStandAlone`), and run `scripts/consumer-smoke.sh` after changing `assets/`, `cmd/ads2css` or `Base`.
- Releases are git tags (`v0.1.0` …) on `main`, made by the owner's request, with the `CHANGELOG.md` entry dated in the same commit. The Go module mirror keeps every tagged version forever, so never move or reuse a tag; fix forward with a new patch version.

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
- The sidebar's collapsed state lives in the `sidebar_state` cookie; the server renders it only through `patterns.SidebarState` (wired in `main`). Tests that call `newMux` directly skip it. Cookies on `localhost` are shared across ports, so a state set on :8090 shows up on :7331 and :8092 too.
- Body is boosted: a link to a file download needs the `download` attribute (htmx then leaves it alone), or the file is swapped into the page.
- Build query strings with `url.Values`, never `"?a=" + x`.
- Filter forms trigger on `filter-change` (`data-filter-form`), never plain `change`: an unnamed combobox input fires `change` on blur, and the refetch replaced the Details sheet as it opened. Lazy fragments in filtered results carry the filter values in their URL.
- A toast showing during a boosted navigation froze the tab (the toast's observers loop once its viewport is gone); `toast-guard.js` prevents it. If you ever see the tab hang right after a click on a link, check that it is loaded.

**Verifying UI**
- Check toasts, dialogs and dark mode in Chrome; HTTP only covers server responses. Check `list_connected_browsers` if the extension reports not connected.
- Dark Reader (if installed) darkens pages even without `.dark`. Judge themes by reading CSS variables, not screenshot colours.
- Dark tokens and the `dark:` variant both follow `.dark` on `<html>` (not `prefers-color-scheme`, not `data-theme`); `theme.js` sets the class from the OS setting and the saved toggle.

## Rules for UI work in this repo

- **Design skills run in polish mode.** When a design skill (unslop, frontend-design) reviews or edits this repo, `docs/design-system.md` and `tokens/tokens.json` override the skill's tokens and checklist. Fix only real signs the skill finds; report anything that conflicts with an approved rule (serif titles at 600, 10px card radius, Lucide stroke 2, navy dark surfaces, sentences with periods) as a suggestion for the design-system owner instead of changing it. A rule changes in `docs/design-system.md` first, then in the code. The repo's copy of unslop (`.claude/skills/unslop/SKILL.md`) carries the same rule in its step 3.

- Use shadcn-templ components first; customise via `Class` and `Attributes` (htmx attributes go in `Attributes`). Don't fork components.
- Colours only through tokens (`bg-primary`, `text-muted-foreground`, `text-positive`, `bg-ink`…). No hex values or Tailwind palette colours (`blue-600`) in templates.
- One `primary` button per view. Destructive actions go through `alertdialog`, never `hx-confirm`.
- Numbers: Arion format in every UI language (`.` thousands, `,` decimal; use `patterns.Int/Dec/Signed/Pct`), `font-mono tabular-nums` and right-aligned in tables (KPI values in `patterns.Stat` are sans `tabular-nums`); gains `text-positive` with `+`, losses `text-destructive` with `−`. Words and units follow the page language ("41 min" on English pages).
- One name per thing: sidebar label = breadcrumb = page title; an action keeps its name from button to toast. No `·`-joined metadata (the design-rules test rejects it); write a sentence.
- Row actions go in the ⋯ menu, destructive ones last after a separator; a menu item opens its dialog with `data-templ-controls` (see `opens` in `ui/pages/showcase.templ`). Page-level actions don't act on one row.
- State pills: `patterns.StatusBadge`, never ad-hoc badge variants. Page chrome: `patterns.Shell`, `PageHeader`, `Stat`, `DataCard` before hand-rolling.
- Type: use the scale utilities (`font-display text-display-sm` for the page title, `text-body-sm`, `text-code`…), never `text-[40px]` or Tailwind's `text-sm`/`text-xs` (the design-rules test rejects both). Everything but titles is sans.
- htmx 4: inheritance is explicit (`hx-x:inherited`), 4xx/5xx swap by default (use `hx-status:422`), use `<hx-partial>` for multi-target responses (toasts go to `#toasts`).
- After changing `tokens/tokens.json`: `task theme`.
