# Patterns

Compositions of shadcn-templ components that every app needs, in `github.com/tilraunastofan/ads2/patterns`. Apps import it from the ads2 module (`go get github.com/tilraunastofan/ads2`, see `docs/setup.md`); they never copy it, so `go get -u` brings every fix. They add no styling of their own: only token classes and component props.

| Pattern | Use |
| --- | --- |
| `Shell(ShellProps)` | The internal-tool document and layout: sidebar with navigation (`Groups []NavGroup`, or `GroupLabel` + `Nav` for one group), theme toggle, one 1200px column. A `NavItem` may have `Children`, listed under it and hidden when the sidebar collapses. The active item is the one whose `Href` is the longest match for `CurrentPath`, the path itself or a segment prefix of it (`/apps` on `/apps/payments`; `/` only on itself); a parent stays marked while one of its children is active. `Footer` items sit at the bottom of the sidebar, outside the scrolling groups (About, Account), and are marked active the same way; `FooterNote` is one small muted line under them, such as the version, hidden when the sidebar collapses. `Base(title, lang)` is the bare document. |
| `Banner`, `WithBanner(ctx, func() *Banner)` | A notice about the whole tool ("Deploys are paused."), which `Shell` shows above every page's content: warning, or `Destructive` for a failure, with an optional link. Set it in middleware; the func runs once per page that renders a `Shell`. Errors that belong to one page go to `#page-alerts` instead. |
| `ReportShell(ReportShellProps)` / `ReportCover(ReportCoverProps)` | Public and report pages: top bar (mark, product name, theme toggle), optional full-width `sand` cover with the one `display-lg` title, and a centred 1040px column. Words follow `Lang`; Icelandic finance views set `Lang: "is"`. Text on the cover is `foreground` and `slate`, never `muted-foreground` (4.3:1 on sand). |
| `PageHeader(PageHeaderProps)` | Breadcrumb, the one `font-display text-display-sm` title, a status line, and up to two actions as children. |
| `Stat(StatProps)` / `StatRow()` | KPI card: metric name, the value at `heading-1`, then the change with sign, arrow and colour and `Since` ("vs last week"). `Good` and `Up` are independent. |
| `StatusBadge(Status, label)` | A state as a pill with icon and word. `StatusPositive` (healthy), `StatusWarning` (progressing), `StatusDestructive` (degraded, failed), `StatusNeutral`. |
| `DataCard(DataCardProps)` | A card holding a full-width table, with optional description and header action (search, button). First and last columns get the card inset: no `pl-4` / `pr-4` on cells. |
| `SidebarState(http.Handler)` | Middleware: reads `sidebar.js`'s `sidebar_state` cookie so `Shell` renders the sidebar collapsed or expanded as the user left it. Wrap your mux in it, or the sidebar re-expands on every navigation. |
| `ThemeToggle(lang...)` | Ghost icon button wired to `assets/js/theme.js`, named "Dark theme" / "Dökkt þema" with `aria-pressed`. Already in `Shell`. |
| `SkipLink(lang)` | "Skip to content" link, visible on focus, to `#main` (`tabindex="-1"`). Already in `Shell`, where `#main` is a `<div>` inside the sidebar inset, the page's one `<main>`; on `Base` pages give your own `<main>` that id. |
| `ChartNumber(v, decimals)`, `ChartTooltipRow(name, color, value)` | Arion numbers in `components/chart`, whose tooltip defaults to en-US: `ChartNumber` for ticks and labels (takes the chart's `any`), `ChartTooltipRow` for the tooltip `Formatter`. See `docs/components/Chart.md`. |
| `Int`, `Dec`, `Signed`, `Pct`, `SignedPct` | Arion number format: `1.284.500`, `2,3`, `+9,5`, `2,3%`, real minus `−`. Any finite value formats exactly, however large; NaN and ±Inf (such as a change from zero) give `patterns.NotAvailable` (`–`), never `0`. |
| `Date`, `DateTime`, `DateRange`, `ISODate`, `ParseISODate` | Dates in the page language: `23. október 2026` / `23 Oct 2026`, `… kl. 14:32` / `…, 14:32`, ranges `20.–24. október 2026` / `20–24 Oct 2026`; ISO `2026-10-23` on the wire. |

```templ
@patterns.PageHeader(patterns.PageHeaderProps{
    Crumbs: []patterns.Crumb{{Label: "Platform", Href: "/"}, {Label: "Releases"}},
    Title:  "Releases",
    Status: "42 applications on ocp-prod-1, synced 1 minute ago.",
}) {
    @button.Button(button.Props{Variant: button.VariantOutline}) { @icon.Download() Export }
    @button.Button() { @icon.Plus() New application }
}
@patterns.StatRow() {
    @patterns.Stat(patterns.StatProps{Label: "Change failure rate", Value: patterns.Pct(2.3, 1),
        Delta: patterns.Signed(-0.8, 1) + " pp", Good: true, Up: false})
}
@patterns.DataCard(patterns.DataCardProps{Title: "Applications", Action: searchBox}) {
    @table.Table() { … }
}
```

## Rules

- `Shell` renders its own words (skip link, sidebar and theme toggle names) in `ShellProps.Lang`: set `Lang: "is"` on Icelandic pages.
- Map your own states onto `Status` in one helper (`Healthy` → positive, `Progressing` → warning, `Degraded` → destructive); don't pick badge variants per page.
- A new app page starts from these and the **Dashboard** composition (`docs/components/Dashboard.md`). If two apps hand-roll the same block, it belongs here.
- `reference`: `ui/pages/showcase.templ`, `metrics.templ` and `settings.templ` use all of them; `ui/pages/report.templ` (`/uppgjor`) is the Icelandic report on `ReportShell`.
