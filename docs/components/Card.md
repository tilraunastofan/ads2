# Card

A bounded group of related content: `card.Card`, `Header`, `Title`, `Description`, `Action`, `Content`, `Footer` from `components/card`.

**Consumer provides:** title, optional description and action (badge, button, menu), content, optional footer; `Size` (`SizeDefault`, `SizeSm`).

```templ
@card.Card(card.Props{ID: "app-" + app.Name, Attributes: templ.Attributes{"hx-get": "/apps/" + app.Name + "/card", "hx-trigger": "every 15s", "hx-swap": "outerHTML"}}) {
    @card.Header() {
        @card.Title() { { app.Name } }
        @card.Description() { Synced { app.SyncedAgo } }
        @card.Action() { @StatusBadge(app.Health) }
    }
    @card.Content() { … }
}
```

- `card` surface on the `background` canvas; separation comes from the 10% foreground ring, not shadows.
- KPI cards: description = metric name, title = value (`tabular-nums`, Icelandic number format), action = trend badge with sign and arrow.
- Poll or SSE-refresh a single card by giving it an `ID` and swapping `outerHTML`.
- Do not nest cards; use `separator` or `item` inside.
