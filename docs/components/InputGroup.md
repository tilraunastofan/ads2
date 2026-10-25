# InputGroup

An input with addons (icon, text, button): `inputgroup.InputGroup`, `inputgroup.Input`, `inputgroup.Addon`, `inputgroup.Button`, `inputgroup.Text` from `components/inputgroup`.

**Consumer provides:** one `inputgroup.Input` or `inputgroup.Textarea`, addons with `Align` (`inline-start`, `inline-end`, `block-start`, `block-end`).

```templ
@inputgroup.InputGroup() {
    @inputgroup.Input(inputgroup.InputProps{Type: "search", Name: "q", Placeholder: "Search services",
        Attributes: templ.Attributes{"hx-get": "/services", "hx-trigger": "input changed delay:300ms, search", "hx-target": "#results"}})
    @inputgroup.Addon(inputgroup.AddonProps{Align: inputgroup.AlignInlineStart}) { @icon.Search() }
}
```

- Search fields: leading search icon, live results via `hx-trigger="input changed delay:300ms"`.
- Currency: trailing `kr.` text, `tabular-nums`, Icelandic grouping (`250.000`).
