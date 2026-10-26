# Combobox

A searchable select for long lists: `combobox.Combobox`, `Input`, `Content`, `List`, `Item{Value, Label}`, `Empty`, `Group`, `Label`, `Separator`, and `Chips` for multiple values, from `components/combobox` (`shadcn-templ add combobox`).

**Consumer provides:** `Name` (submits a hidden input), `DefaultValue` or `DefaultValues` with `Multiple`, the items, a placeholder, and the empty text.

```templ
@combobox.Combobox(combobox.Props{Name: "app", DefaultValue: current}) {
    @combobox.Input(combobox.InputProps{Placeholder: "All applications", ShowClear: true})
    @combobox.Content() {
        @combobox.Empty() { No matching application }
        @combobox.List() {
            for _, a := range apps {
                @combobox.Item(combobox.ItemProps{Value: a.Name}) { { a.Name } }
            }
        }
    }
}
```

- Use for more than about ten options or when people type to find (applications, namespaces, cost centres). Short lists: `nativeselect`.
- Filtering is client-side over the rendered items. For thousands of items render a server-side search (`inputgroup` + `hx-get`) instead.
- Changing the value fires `change` on the hidden input, so a surrounding `hx-trigger="change"` form refreshes the page region.
- Solid `popover` menu, as everywhere in this system.
