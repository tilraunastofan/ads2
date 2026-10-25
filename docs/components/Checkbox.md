# Checkbox

A boolean or multi-select option: `checkbox.Checkbox` from `components/checkbox` (Base UI checkbox: `role="checkbox"`, `data-checked`, hidden native input for forms).

**Consumer provides:** `ID`, `Name`, `Value`, `DefaultChecked` (or `Checked *bool`), `Disabled`, `Indeterminate`, and a label.

```templ
@field.Field(field.Props{Orientation: field.OrientationHorizontal}) {
    @checkbox.Checkbox(checkbox.Props{ID: "migrate", Name: "migrate", DefaultChecked: true})
    @field.Label(field.LabelProps{For: "migrate"}) { Run database migrations }
}
```

- Checked fill is `primary`; unchecked border is `input`.
- For an immediate on/off setting use `switch`; checkboxes are for choices submitted with a form.
- Use `Indeterminate` on a "select all" header checkbox in tables.
