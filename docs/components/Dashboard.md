# Dashboard

A composition, not a component: the page pattern for internal platform and finance dashboards built only from shadcn-templ parts (`breadcrumb`, `button`, `card`, `badge`, `table`, `dropdownmenu`).

**Consumer provides:** the page title, a one-line status sentence, two header actions at most, three or four KPI cards, and the main table.

```templ
templ ReleasesPage(vm ReleasesVM) {
    @layouts.App("Releases") {
        @PageHeader(vm.Crumbs, "Releases", vm.Status) {
            @button.Button(button.Props{Variant: button.VariantOutline, Href: "/releases.csv"}) { @icon.Download() Export }
            @button.Button(button.Props{Href: "/apps/new"}) { @icon.Plus() New application }
        }
        <div class="grid gap-4 sm:grid-cols-3">
            for _, k := range vm.KPIs {
                @KPICard(k)
            }
        </div>
        @card.Card() {
            @card.Header(card.HeaderProps{Class: "border-b"}) { @card.Title() { Applications } }
            <div id="apps" hx-get="/releases/table" hx-trigger="every 15s" hx-swap="outerHTML">
                @AppsTable(vm.Apps)
            </div>
        }
    }
}
```

- Page title in the `display` family (`display-sm`, 40px serif) — the one Arion signature on every page. Everything else is the sans UI.
- Order: breadcrumb → title + status line + actions → KPI row → main table. 24px (`space-6`) between blocks.
- KPI deltas carry sign, arrow and colour (`positive` / `destructive`); "good" direction depends on the metric (a falling failure rate is green).
- Only the table polls; the shell never re-renders.
