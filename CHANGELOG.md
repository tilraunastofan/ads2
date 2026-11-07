# Changelog

What changed in each release of `github.com/tilraunastofan/ads2`, and what an app has to do to take it. Apps update with `go get -u github.com/tilraunastofan/ads2` (see **Updating later** in `docs/migrate-to-module.md`).

Versions follow [semver](https://semver.org). Before `v1.0.0` a minor version (`v0.2.0`) may break the API, and its entry says how to update; a patch (`v0.1.1`) never does.

Each entry has **Breaking** (an app must change code), **Added** and **Fixed**. Keep entries short and concrete: name the function, file or token, and the one thing an app does about it.

## Unreleased

**Added**

- `patterns.Base` loads the htmx 4 extensions `hx-preload`, `hx-history-cache` and `hx-ptag` (vendored in `assets/js/ext/`) and an `htmx-config` meta that turns on preload for boosted links on `mouseover`. Boosted links are fetched when the pointer reaches them, and Back/Forward restores pages from `sessionStorage` (see **Extensions the shell loads** in `docs/htmx-patterns.md`). A GET behind a boosted link may now run on hover, so it must have no side effects beyond audit logging. An app that loads these extensions or sets `htmx-config` itself drops its own copy.

- `ShellProps.Background` (`Background{Light, Dark}`) and `ArcticBackground`: a decorative image behind the page content, below the top bar, that follows the theme (`Light` in light mode, `Dark` under `.dark`; one set shows in both). `ArcticBackground` uses the new `assets/img/arctic-background.svg` and `arctic-background-light.svg`, the Arion arctic landscape without the star. For calm pages (home, help, about), not data pages. Unset, the shell is unchanged.

- `WithHeaderActions(ctx, func() templ.Component)` and `HeaderActionsFrom(ctx)`: an app puts a component (a notifications bell) in `Shell`'s top bar, before the theme toggle. Set it in middleware, like `WithBanner`; the func runs once per render. Unset, the header is unchanged.

**Fixed**

- `chart.Line`: the legend swatch and the default tooltip indicator ignored `LineProps.Stroke` and read `--color-<key>` from `Config.Color`, so under a strict CSP (no `Config.Color`) they were blank. They now take the stroke, like Recharts' `getLegendItemColor` and as `docs/components/Chart.md` already said; a `Stroke` of `""` or `"none"` still falls back to the series colour.

Nothing to change in apps: `go get -u`. An app that set `--color-<key>` on the container's `Attributes` to colour a line's swatch can drop that.

## v0.1.8 (4 Nov 2026)

**Added**

- `ShellProps.Footer` and `ShellProps.FooterNote`: items pinned to the bottom of the sidebar, below the scrolling groups, and a small muted line under them (for example the app's version). Footer items look and mark the active page like the other nav items; the note hides when the sidebar collapses to icons. Both are empty by default and then nothing renders.

Nothing to change in apps: `go get -u`, then rebuild the CSS. To use it, set `Footer` and/or `FooterNote` in your `ShellProps`.

## v0.1.7 (4 Nov 2026)

**Fixed**

- `filter-form.js`: a live-search input (`hx-trigger="input changed delay:300ms, search"`) inside a `form[data-filter-form]` sent its request while typing, and the native `change` on blur then made the form refetch the same request and, with `hx-push-url`, push a duplicate history entry (Back needed two presses). The script now refreshes its snapshot on `htmx:before:request` from any element in the form.

Nothing to change in apps: `go get -u`.

## v0.1.6 (4 Nov 2026)

Repo tooling only (`task clean`, and `task check` rebuilds the minified script bundle). Nothing in the module changed; apps need not update.

## v0.1.5 (4 Nov 2026)

**Added**

- `components/hovercard` (shadcn-templ's `hover-card`, Base UI's PreviewCard): read-only detail on hover or focus, portaled to `<body>` so table cells and scroll containers don't clip or squeeze it. See `docs/components/HoverCard.md`; the showcase's `/apps` version column uses it. The component script bundle is now `assets/js/shadcn-templ-85e9a7167c0ea409.js`; `components.Scripts()` links it.

Nothing to change in apps: `go get -u`, then rebuild the CSS (`go tool ads2css` picks up the new component's classes).

## v0.1.4 (4 Nov 2026)

**Fixed**

- `Shell`: a wide table or code block made the whole page scroll sideways at tablet width (768px, sidebar open) instead of scrolling inside its card. The sidebar inset now has `min-w-0`.
- `PageHeader`: a long title without spaces (`IsIT.AML.Services`) ran off a phone screen; it now wraps (`break-words`, and `min-w-0` on its column).

Nothing to change in apps: `go get -u`, then rebuild the CSS.

## v0.1.3 (3 Nov 2026)

**Fixed**

- The component script bundle (`assets/js/shadcn-templ-<hash>.js`) is in the module. It was gitignored, so v0.1.0 to v0.1.2 shipped without it: `components.Scripts()` linked a file that `assets.Handler` answered with a 404, and no component script ran (sidebar toggle, dialogs, menus, toasts). Apps update with `go get -u`; nothing else changes.
- `scripts/consumer-smoke.sh` builds the consumer against the repo's tracked and new files, as the published module has them, instead of the working tree, so an ignored file the module needs fails it. `task check` also fails on an uncommitted bundle.

## v0.1.2 (3 Nov 2026)

**Added**

- `assets/js/copy.js`, loaded by `patterns.Base`: a button with `data-copy-target="<id>"` copies that element's text and shows "Copied" / "Afritað" in its `[data-label]` (see **Copy button** in `docs/htmx-patterns.md`).
- `assets/js/submit-once.js`, loaded by `patterns.Base`: `form[data-disable-on-submit]` disables its submit buttons while it submits, with `data-busy-label` in the button's `[data-label]` (see **Submit once**). An app that had its own `copy.js` or `submit-once.js` deletes it and stops loading it, or it runs twice.
- Token `orange` (`#b34700` / `#ff9a4d`) and `orange-foreground`: the step between `destructive` and `warning` in a graded scale (Urgent / Critical / Serious), never a status of its own. Contrast-checked like `warning`. Replace `orange-500`/`orange-600` palette classes with `text-orange` / `bg-orange/10`.

## v0.1.1 (3 Nov 2026)

**Added**

- `NavItem.Children`: sub-items listed under a sidebar item, always expanded, hidden when the sidebar collapses to icons.
- `patterns.Banner`, `WithBanner` and `BannerFrom`: a page-wide notice (warning, or destructive) that `Shell` shows above the page's content, with an optional link. Set it in middleware. An app that rendered its own banner above `Shell`'s children moves it here.

**Changed**

- `Shell` marks the item whose `Href` is the longest match for `CurrentPath`, the path itself or a segment prefix of it, instead of only an exact match: `/apps/payments` now marks `/apps`. `/` still matches only itself. Nothing to change unless a page relied on no item being marked below an item's path.

## v0.1.0 (3 Nov 2026)

The first release, and the first one apps import as a Go module instead of copying. Apps that copied `components/`, `patterns/` and the assets: follow `docs/migrate-to-module.md`.

**Added**

- Apps import the design system: `github.com/tilraunastofan/ads2/patterns`, `/components/<name>`, `/assets` (see **How apps use ads2** in `docs/setup.md`).
- `assets.Handler(fsys)` serves the theme, fonts, htmx, helper scripts and component bundle with ETags, gzip and a year-long cache for fonts and the hashed bundle. `assets.Overlay(app)` puts the app's own files (its `css/output.css`) on top. `assets.AcceptsGzip` and `assets.CacheControl` are the rules it uses.
- `cmd/ads2css` (`go tool ads2css`) writes the CSS file that connects an app's Tailwind build to the module: the module's `assets/css/ads2.css` and `@source` lines for its components and patterns.
- `patterns.ChartNumber` and `patterns.ChartTooltipRow`: Arion numbers in `components/chart`, whose tooltip defaults to en-US. Pass `ChartTooltipRow` as the tooltip `Formatter` in every chart (`docs/components/Chart.md`).
- `patterns.NotAvailable` (`–`), what `Dec`, `Signed`, `Pct` and `SignedPct` return for NaN and ±Inf.
- `assets/js/focus-after-swap.js`, loaded by `patterns.Base`: focus moves to the new page's `<h1>` after a boosted navigation, and to the first invalid field of a form re-rendered with a 422.

**Changed** (for apps migrating from an older copy)

- `patterns.Base` loads every helper script, `datepicker.js` included. Remove `<script src="/assets/js/datepicker.js">` from pages: a script inside a page runs again on every boosted visit.
- `PageHeader`'s `<h1>` has `tabindex="-1"` (the focus target after navigation).
- `Shell`'s skip-link target `#main` is a `<div>` inside the sidebar inset, which is the page's one `<main>`.
- `aria-invalid` must be the string `"true"`, never a Go bool (a bool renders a bare attribute, which ARIA reads as false). The design-rules test rejects bools.
- The type scale replaces Tailwind's text sizes in patterns (`text-body-sm`, `text-caption`); do the same in app templates.
- Bars in charts: `Radius: []float64{4, 4, 0, 0}` (square on the baseline).
- Token: dark `chart-5` is `#8ea0ab`, a slate (was near-white `#e7e0d8`).

**Fixed**

- `Dec` and `Signed` format any finite value exactly (1e20 was clamped to `math.MaxInt64`).
- `Shell` renders a `NavItem` without an `Icon` instead of panicking.
- `filter-form.js` no longer refetches on the first blur of a filter form.
- `AcceptsGzip` refuses `gzip;q=0.0` (only `q=0` was refused).

Built on shadcn-templ `v2.0.0`, htmx `4.0.1`, templ `v0.3.1100`. The vendored toast still has the freeze reported as [axadrn/shadcn-templ#656](https://github.com/axadrn/shadcn-templ/issues/656); `assets/js/toast-guard.js`, loaded by `Base`, prevents it.
