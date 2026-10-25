# Toggle

A two-state button: `toggle.Toggle` from `components/toggle`, with `aria-pressed` and `data-pressed`.

**Consumer provides:** label or icon children, `DefaultPressed` (or a `Pressed` pointer when the server owns the state), `Variant` (`VariantDefault`, `VariantOutline`), `Size`, `Value`.

```templ
@toggle.Toggle(toggle.Props{
    Variant: toggle.VariantOutline,
    DefaultPressed: filters.OnlyFailing,
    Attributes: templ.Attributes{"hx-get": "/pipelines", "hx-include": "closest form", "hx-trigger": "click"},
}) { Only failing }
```

- Use for view filters and formatting, not for settings that persist (use `switch`).
- Icon-only toggles need an `aria-label`.
