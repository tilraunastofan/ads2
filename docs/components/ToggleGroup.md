# ToggleGroup

A set of toggles with single or multiple selection: `togglegroup.ToggleGroup` + `togglegroup.Item` from `components/togglegroup`.

**Consumer provides:** items with `Value`, `DefaultValue []string`, `Multiple`, `Variant`, `Size`, `Spacing` (0 joins them like a segmented control).

```templ
{{ zero := 0 }}
@togglegroup.ToggleGroup(togglegroup.Props{DefaultValue: []string{"dev"}, Variant: togglegroup.VariantOutline, Spacing: &zero}) {
    @togglegroup.Item(togglegroup.ItemProps{Value: "dev"}) { Dev }
    @togglegroup.Item(togglegroup.ItemProps{Value: "staging"}) { Staging }
    @togglegroup.Item(togglegroup.ItemProps{Value: "prod"}) { Prod }
}
```

- Environment and time-range switchers; three to five short options.
- Pair with `hx-get` on each item when the selection reloads a region.
