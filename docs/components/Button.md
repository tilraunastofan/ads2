# Button

The action control: `button.Button` from `github.com/axadrn/shadcn-templ/v2/components/button`, rendered as `<button>` or, with `Href`, as `<a role="button">`.

**Consumer provides:** children (label, optional leading icon from `components/icon`), `Variant`, `Size`, `Type`, and htmx attributes through `Attributes`.

```templ
@button.Button(button.Props{
    Type:       button.TypeSubmit,
    Attributes: templ.Attributes{"hx-post": "/releases/42/promote", "hx-disable": "this"},
}) {
    Promote to prod
}
```

- One `VariantDefault` (Arion blue `primary`) per view: the action the screen exists for.
- `VariantSecondary` for the second action, `VariantOutline` for tool actions (export, filter), `VariantGhost` for cancel and toolbar icons.
- `VariantDestructive` only for irreversible actions, and always behind `alertdialog`.
- `VariantLink` inside prose and tables, never as a primary action.
- Icon-only buttons (`SizeIcon*`) need `templ.Attributes{"aria-label": "…"}`.
- Arion's dark editorial CTA: `Class: "bg-ink text-ink-foreground hover:bg-ink-hover"` on `VariantDefault`. One per page, never next to a blue primary.
- While an htmx request runs, show a spinner and disable: `hx-disable="this"` plus an `icon.LoaderCircle(icon.Props{Class: "animate-spin"})` swapped in, or a `.htmx-request` style.
- Labels are verbs in sentence case: "Promote to prod", not "PROMOTE" or "Click here".
