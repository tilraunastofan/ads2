# Dashboard

A composition, not a component: the page pattern for internal platform and finance dashboards, built from the `patterns` package (`PageHeader`, `Stat`, `DataCard`, `StatusBadge`) and shadcn-templ parts.

**Consumer provides:** the page title, a one-line status sentence, two header actions at most, three or four KPI cards, and the main table.

```templ
templ ReleasesPage(vm ReleasesVM) {
    @layouts.App("Releases", "/releases") {
        @patterns.PageHeader(patterns.PageHeaderProps{Crumbs: vm.Crumbs, Title: "Releases", Status: vm.Status}) {
            @button.Button(button.Props{Variant: button.VariantOutline, Href: "/releases.csv"}) { @icon.Download() Export }
            @button.Button(button.Props{Href: "/apps/new"}) { @icon.Plus() New application }
        }
        @patterns.StatRow() {
            for _, k := range vm.KPIs {
                @patterns.Stat(k)
            }
        }
        @patterns.DataCard(patterns.DataCardProps{ID: "apps", Title: "Applications"}) {
            <div hx-get="/releases/table" hx-trigger="every 15s" hx-swap="outerHTML">
                @AppsTable(vm.Apps)
            </div>
        }
    }
}
```

- Page title in the `display` family (`font-display text-display-sm`, 40px serif) — the one Arion signature on every page. `PageHeader` does this. Everything else is the sans UI.
- Order: breadcrumb → title + status line + actions → KPI row → main table. 24px (`space-6`) between blocks.
- KPI deltas carry sign, arrow and colour (`positive` / `destructive`); "good" direction depends on the metric (a falling failure rate is green).
- Only the table polls; the shell never re-renders.
