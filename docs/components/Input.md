# Input

Single-line text entry: `input.Input` from `components/input`.

**Consumer provides:** `Name`, `Type`, `Placeholder`, `Value`, `Required`/`Disabled`/`ReadOnly`, and an `ID` that a `field.Label` points at. Wrap it in `field.Field`.

```templ
@input.Input(input.Props{
    ID: "image", Name: "image", Placeholder: "registry.internal/team/app:tag",
    Attributes: templ.Attributes{"hx-post": "/validate/image", "hx-trigger": "change, keyup delay:400ms", "hx-target": "#image-error"},
})
```

- Always labelled: use `field.Label`, never a placeholder as the label.
- Mark errors with `aria-invalid="true"` (draws the `destructive` border) and show `field.Error` text below.
- Control borders use `input` (Arion's 50% navy).
- Use `inputgroup` for prefixes, suffixes and inline buttons.
