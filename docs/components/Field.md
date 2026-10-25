# Field

Layout and semantics for one form control with label, description and error: `field.Field`, `field.Label`, `field.Description`, `field.Error`, `field.Group`, `field.Set`, `field.Legend`, `field.Content` from `components/field`.

**Consumer provides:** the control (input, select, checkbox…), its label text, optional description, error text when invalid, `Orientation` (`vertical`, `horizontal`, `responsive`).

```templ
{{ attrs := templ.Attributes{} }}
if err != "" {
    {{ attrs["data-invalid"] = "true" }}
}
@field.Field(field.Props{Attributes: attrs}) {
    @field.Label(field.LabelProps{For: "replicas"}) { Replicas }
    @input.Input(input.Props{ID: "replicas", Name: "replicas", Type: "number", Value: v,
        Attributes: templ.Attributes{"aria-invalid": err != ""}})
    if err != "" {
        @field.Error() { { err } }
    }
}
```

- Every form is `field.Group` > `field.Field`s; group related fields in `field.Set` + `field.Legend`.
- Validate on the server and re-render the whole `field.Field` with htmx (`hx-post` on the form, `hx-select` or `<hx-partial>` per field). Return 422 and use `hx-status:422` to keep the swap on the form.
- `field.Description` uses `muted-foreground`; `field.Error` uses `destructive` and is announced (`role="alert"`).
- Horizontal orientation for checkbox and switch rows.
