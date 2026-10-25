# Slider

Pick a number or range by dragging: `slider.Slider` from `components/slider`.

**Consumer provides:** `Name`, `DefaultValue` (number or `[]float64` for a range), `Min`, `Max`, `Step`, `Orientation`, plus a visible label and value readout.

```templ
@slider.Slider(slider.Props{Name: "canary", DefaultValue: 10, Min: 0, Max: 100, Step: 5})
```

- Only where an approximate value is fine (canary weight, sampling rate). For exact amounts use `input`.
- Always show the current value as text next to it.
