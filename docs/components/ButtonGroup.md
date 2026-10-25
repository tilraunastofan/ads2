# ButtonGroup

Joins related buttons into one control: `buttongroup.ButtonGroup` from `components/buttongroup`.

**Consumer provides:** two or more `button.Button` children of the same `Variant`; optional `buttongroup.Text` and `buttongroup.Separator`; `Orientation` (horizontal by default).

```templ
@buttongroup.ButtonGroup() {
    @button.Button(button.Props{Variant: button.VariantOutline, Attributes: templ.Attributes{"hx-get": "/metrics?range=day", "hx-target": "#chart"}}) { Day }
    @button.Button(button.Props{Variant: button.VariantOutline, Attributes: templ.Attributes{"hx-get": "/metrics?range=week", "hx-target": "#chart"}}) { Week }
}
```

- Use for split buttons (action + `dropdownmenu` trigger) and segmented filters that each fire a request.
- For a persistent on/off selection use `togglegroup` instead: it carries `aria-pressed` state.
- Never mix variants inside one group.
