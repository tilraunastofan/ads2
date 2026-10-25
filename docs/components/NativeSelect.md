# NativeSelect

A styled native `<select>`: `nativeselect.NativeSelect`, `nativeselect.Option`, `nativeselect.OptGroup` from `components/nativeselect`.

**Consumer provides:** `Name`, `Size`, options with `Value` and `Selected`.

```templ
@nativeselect.NativeSelect(nativeselect.Props{Name: "namespace",
    Attributes: templ.Attributes{"hx-get": "/workloads", "hx-target": "#workloads", "hx-trigger": "change"}}) {
    for _, ns := range namespaces {
        @nativeselect.Option(nativeselect.OptionProps{Value: ns, Selected: ns == current}) { { ns } }
    }
}
```

- Default choice for server-rendered forms: zero JS, works with `hx-include` and plain form posts.
- Use `select` only when you need rich items (icons, descriptions, groups with custom rendering).
