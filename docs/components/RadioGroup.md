# RadioGroup

One choice from a short list: `radiogroup.RadioGroup` + `radiogroup.Item` from `components/radiogroup`.

**Consumer provides:** `Name`, `DefaultValue` (or `Value *string`), items with `Value` and an `ID` for their labels.

```templ
@radiogroup.RadioGroup(radiogroup.Props{Name: "strategy", DefaultValue: "rolling"}) {
    @field.Field(field.Props{Orientation: field.OrientationHorizontal}) {
        @radiogroup.Item(radiogroup.ItemProps{ID: "s-rolling", Value: "rolling"})
        @field.Label(field.LabelProps{For: "s-rolling"}) { Rolling update }
    }
}
```

- Two to five options that need to be seen side by side; more than five → `select`.
- Always pre-select the safe default.
