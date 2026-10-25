# Alert

An inline, persistent message about the page or a section: `alert.Alert{Variant}`, `Title`, `Description`, `Action` from `components/alert`.

```templ
@alert.Alert(alert.Props{Variant: alert.VariantDestructive}) {
    @icon.CircleAlert()
    @alert.Title() { Sync failed }
    @alert.Description() { { err.Error() } }
    @alert.Action() { @button.Button(button.Props{Variant: button.VariantOutline, Size: button.SizeXs, Attributes: templ.Attributes{"hx-post": retryURL}}) { Retry } }
}
```

- Lead with an icon (`info`, `circle-alert`), a short title, one or two sentences.
- `VariantDestructive` for failures that block the user; transient results go in a `toast`. Success and warning notes use `VariantDefault` with a leading `icon.CircleCheck()` / `icon.TriangleAlert()`, tinted with `Class: "text-positive"` / `"text-warning"` on the icon only.
- Return an alert from an htmx error response with `hx-status:5xx="target:#page-alerts"` (htmx 4 swaps 4xx/5xx by default).
