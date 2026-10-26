# PageHeader

The top of every page: optional breadcrumb, the one serif display title, a status line, and up to two actions. `patterns.PageHeader(PageHeaderProps)` from `github.com/tilraunastofan/ads2/patterns` (copy the package into your app like `components/`).

**Consumer provides:** `Crumbs []Crumb{Label, Href}` (last one is the current page, `Href` empty), `Title`, `Status` (one sentence), and the actions as children.

```templ
@patterns.PageHeader(patterns.PageHeaderProps{
    Crumbs: []patterns.Crumb{{Label: "Platform", Href: "/"}, {Label: "Releases"}},
    Title:  "Releases",
    Status: "ocp-prod-1 · 42 applications · synced 1 minute ago",
}) {
    @button.Button(button.Props{Variant: button.VariantOutline, Href: "/apps.csv"}) { @icon.Download() Export }
    @button.Button() { @icon.Plus() New application }
}
```

- Exactly one per page, first in `main`. The title is `font-display text-display-sm` (Source Serif 4, 40px): the only serif on an app page.
- At most two actions: the primary last. More actions go in a `dropdownmenu`.
- The status line states freshness or scope in plain words; it is not a subtitle slogan.
