# Tabs

Switch between views of the same object: `tabs.Tabs`, `tabs.List`, `tabs.Trigger`, `tabs.Content` from `components/tabs`.

**Consumer provides:** `DefaultValue`, a `List` (`Variant`: `VariantDefault` pill or `VariantLine` underline) of `Trigger{Value}`, and one `Content{Value}` per trigger.

```templ
@tabs.Tabs(tabs.Props{DefaultValue: "overview"}) {
    @tabs.List() {
        @tabs.Trigger(tabs.TriggerProps{Value: "overview"}) { Overview }
        @tabs.Trigger(tabs.TriggerProps{Value: "logs"}) { Logs }
    }
    @tabs.Content(tabs.ContentProps{Value: "overview"}) { @OverviewPanel(app) }
    @tabs.Content(tabs.ContentProps{Value: "logs"}) {
        <div hx-get={ "/apps/" + app.Name + "/logs" } hx-trigger="intersect once">@skeleton.Skeleton(skeleton.Props{Class: "h-48"})</div>
    }
}
```

- `VariantDefault` inside cards and dialogs; `VariantLine` for page-level sections (Arion's report tabs).
- Lazy-load heavy panels with `hx-trigger="intersect once"`.
- Tabs are not navigation between pages: use `sidebar` or links for that.
