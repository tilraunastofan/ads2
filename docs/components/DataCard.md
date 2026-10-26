# DataCard

A card whose content is a table running edge to edge, with a title, optional description and one header action. `patterns.DataCard(DataCardProps)` from `github.com/tilraunastofan/ads2/patterns` (copy the package into your app like `components/`).

**Consumer provides:** `ID` (to swap the whole card), `Title`, `Description`, `Action` (a `templ.Component`: search box, button or menu), and the `table.Table` as children.

```templ
@patterns.DataCard(patterns.DataCardProps{ID: "apps", Title: "Applications", Description: "Argo CD sync and health", Action: AppSearch(q)}) {
    @table.Table() {
        @table.Header() { … }
        @AppRows(apps)   // table.Body with ID "app-rows", the htmx swap target
    }
}
```

- First and last columns get the card's 16px inset automatically: don't add `pl-4` / `pr-4` to cells.
- One table per card. Paging goes under the table in a `card.Footer`.
- Search and filters in `Action` target the body (`hx-target="#app-rows"`), not the card.
