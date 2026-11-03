# Patterns

Compositions of shadcn-templ components that every app needs, in `github.com/tilraunastofan/ads2/patterns`. Copy the package into a project like `components/` and rewrite its import paths to your module. They add no styling of their own: only token classes and component props.

| Pattern | Use |
| --- | --- |
| `Shell(ShellProps)` | The internal-tool document and layout: sidebar with navigation, theme toggle, one 1200px column. `Base(title, lang)` is the bare document. |
| `ReportShell(ReportShellProps)` / `ReportCover(ReportCoverProps)` | Public and report pages: top bar (mark, product name, theme toggle), optional full-width `sand` cover with the one `display-lg` title, and a centred 1040px column. Words follow `Lang`; Icelandic finance views set `Lang: "is"`. Text on the cover is `foreground` and `slate`, never `muted-foreground` (4.3:1 on sand). |
| `PageHeader(PageHeaderProps)` | Breadcrumb, the one `font-display text-display-sm` title, a status line, and up to two actions as children. |
| `Stat(StatProps)` / `StatRow()` | KPI card: metric name, value, trend badge with sign, arrow and colour. `Good` and `Up` are independent. |
| `StatusBadge(Status, label)` | A state as a pill with icon and word. `StatusPositive` (healthy), `StatusWarning` (progressing), `StatusDestructive` (degraded, failed), `StatusNeutral`. |
| `DataCard(DataCardProps)` | A card holding a full-width table, with optional description and header action (search, button). First and last columns get the card inset: no `pl-4` / `pr-4` on cells. |
| `ThemeToggle(lang...)` | Ghost icon button wired to `assets/js/theme.js`, named "Dark theme" / "Dökkt þema" with `aria-pressed`. Already in `Shell`. |
| `SkipLink(lang)` | "Skip to content" link, visible on focus, to `<main id="main" tabindex="-1">`. Already in `Shell`; add it yourself on `Base` pages. |
| `Int`, `Dec`, `Signed`, `Pct`, `SignedPct` | Arion number format: `1.284.500`, `2,3`, `+9,5`, `2,3%`, real minus `−`. |
| `Date`, `DateTime`, `DateRange`, `ISODate`, `ParseISODate` | Dates in the page language: `23. október 2026` / `23 Oct 2026`, `… kl. 14:32` / `…, 14:32`, ranges `20.–24. október 2026` / `20–24 Oct 2026`; ISO `2026-10-23` on the wire. |

```templ
@patterns.PageHeader(patterns.PageHeaderProps{
    Crumbs: []patterns.Crumb{{Label: "Platform", Href: "/"}, {Label: "Releases"}},
    Title:  "Releases",
    Status: "ocp-prod-1 · 42 applications · synced 1 minute ago",
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
