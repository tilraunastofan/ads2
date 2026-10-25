# Tooltip

A short label on hover/focus: `tooltip.Tooltip`, `tooltip.Trigger(ctx)`, `tooltip.Content{Side, SideOffset}` from `components/tooltip`.

```templ
@tooltip.Tooltip() {
    @button.Button(button.Props{Variant: button.VariantOutline, Size: button.SizeIcon, Attributes: tooltip.Trigger(ctx)}) { @icon.Copy() }
    @tooltip.Content() { Copy image tag }
}
```

- Names icon-only buttons and truncated values. Never holds information that exists nowhere else, and never interactive content (use `popover`).
- Inverted colours: `foreground` fill with `background` text.
