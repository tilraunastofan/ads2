# Patterns

Compositions of shadcn-templ components that every app needs, in `github.com/tilraunastofan/ads2/patterns`. Copy the package into a project like `components/` and rewrite its import paths to your module. They add no styling of their own: only token classes and component props.

| Pattern | Use |
| --- | --- |
| `Shell(ShellProps)` | The internal-tool document and layout: sidebar with navigation, theme toggle, one 1200px column. `Base(title, lang)` is the bare document. |
| `PageHeader(PageHeaderProps)` | Breadcrumb, the one `font-display text-display-sm` title, a status line, and up to two actions as children. |
| `Stat(StatProps)` / `StatRow()` | KPI card: metric name, value, trend badge with sign, arrow and colour. `Good` and `Up` are independent. |
| `StatusBadge(Status, label)` | A state as a pill with icon and word. `StatusPositive` (healthy), `StatusWarning` (progressing, degraded), `StatusDestructive` (failed), `StatusNeutral`. |
| `DataCard(DataCardProps)` | A card holding a full-width table, with optional description and header action (search, button). First and last columns get the card inset: no `pl-4` / `pr-4` on cells. |
| `ThemeToggle()` | Ghost icon button wired to `assets/js/theme.js`. Already in `Shell`. |
| `Int`, `Dec`, `Signed`, `Pct`, `SignedPct` | Arion number format: `1.284.500`, `2,3`, `+9,5`, `2,3%`, real minus `−`. |

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

- Map your own states onto `Status` in one helper (`Healthy` → positive, `Progressing` → warning, `Degraded` → destructive); don't pick badge variants per page.
- A new app page starts from these and the **Dashboard** composition (`docs/components/Dashboard.md`). If two apps hand-roll the same block, it belongs here.
- `reference`: `ui/pages/showcase.templ` uses all of them.
