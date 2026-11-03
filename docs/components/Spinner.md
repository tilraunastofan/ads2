# Spinner

Indeterminate loading: `spinner.Spinner` from `components/spinner` (an animated `icon.LoaderCircle`, `role="status"`).

```templ
@button.Button(button.Props{Variant: button.VariantOutline, Attributes: templ.Attributes{
    "hx-post": "/sync", "hx-swap": "none", "hx-indicator": "#sync-spin", "hx-disable": "this",
}}) {
    Sync
}
<span id="sync-spin" class="htmx-indicator">
    @spinner.Spinner()
</span>
```

- Pair with `hx-indicator`; htmx adds `.htmx-request` while the request runs.
- Inside buttons put it before the label and disable the button.
