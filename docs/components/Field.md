# Field

Layout and semantics for one form control with label, description and error: `field.Field`, `field.Label`, `field.Description`, `field.Error`, `field.Group`, `field.Set`, `field.Legend`, `field.Content` from `components/field`.

**Consumer provides:** the control (input, select, checkbox…), its label text, optional description, error text when invalid, `Orientation` (`vertical`, `horizontal`, `responsive`).

```templ
// invalid(err) is data-invalid="true" and fieldInput(id, err, extra) is
// aria-invalid="true" plus aria-describedby="<id>-hint", both only while err != ""
// (ui/pages/showcase.templ). The description and the error share the hint id.
@field.Field(field.Props{Attributes: invalid(err)}) {
    @field.Label(field.LabelProps{For: "replicas"}) { Replicas }
    @input.Input(input.Props{ID: "replicas", Name: "replicas", Type: "number", Value: v,
        Attributes: fieldInput("replicas", err, nil)})
    if err != "" {
        @field.Error(field.ErrorProps{ID: "replicas-hint"}) { { err } }
    } else {
        @field.Description(field.DescriptionProps{ID: "replicas-hint"}) { At least 2 in production. }
    }
}
```

- Every form is `field.Group` > `field.Field`s; group related fields in `field.Set` + `field.Legend`.
- Validate on the server and re-render the whole `field.Field` with htmx (`hx-post` on the form, `hx-select` or `<hx-partial>` per field). Return 422 and use `hx-status:422` to keep the swap on the form.
- `field.Description` uses `muted-foreground`; `field.Error` uses `destructive` and is announced (`role="alert"`).
- On the input: `aria-invalid="true"` as a string while there is an error (a Go bool renders a bare attribute, which ARIA reads as false; the design-rules test rejects it), and `aria-describedby` pointing at one id that the description and the error share, so whichever shows is read with the field. `fieldInput` in `ui/pages/showcase.templ` does both.
- Horizontal orientation for checkbox and switch rows.
