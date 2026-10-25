# Label

Accessible name for a control: `label.Label` from `components/label` (`For` = the control's `ID`).

```templ
@label.Label(label.Props{For: "ns"}) { Namespace }
```

- Inside forms prefer `field.Label`, which also handles disabled and invalid states of its `field.Field`.
- Sentence case, no trailing colon.
