# Select

A custom listbox select (Base UI): `selectcomp.Select`, `Trigger`, `Value`, `Content`, `Group`, `Label`, `Item`, `Separator` from `components/select` (import alias needed: `select` is a Go keyword).

**Consumer provides:** `Name`, `DefaultValue`, items with `Value` and `Label`, a `Placeholder` on `Value`, trigger `Size`.

```templ
@selectcomp.Select(selectcomp.Props{Name: "env", DefaultValue: "prod"}) {
    @selectcomp.Trigger() { @selectcomp.Value(selectcomp.ValueProps{Placeholder: "Choose environment"}) }
    @selectcomp.Content() {
        @selectcomp.Item(selectcomp.ItemProps{Value: "dev", Label: "Development"}) { Development }
        @selectcomp.Item(selectcomp.ItemProps{Value: "prod", Label: "Production"}) { Production }
    }
}
```

- Prefer `nativeselect` unless items need custom content. Use `combobox` for long, searchable lists.
- Menus render solid `popover` (not shadcn's translucent menu option) to keep financial data legible.
- Preview is static: the content is shown open below the trigger.
