# Dialog

A modal task that interrupts the page: `dialog.Dialog`, `dialog.Trigger(ctx)` (attributes for your button), `Content`, `Header`, `Title`, `Description`, `Footer`, `dialog.Close(ctx)` from `components/dialog`.

```templ
@dialog.Dialog(dialog.Props{ID: "scale"}) {
    @button.Button(button.Props{Variant: button.VariantOutline, Attributes: dialog.Trigger(ctx)}) { Scale }
    @dialog.Content() {
        <form hx-post={ scaleURL } hx-target="this" hx-swap="outerHTML">
            @dialog.Header() { @dialog.Title() { Scale deployment } @dialog.Description() { Changes apply through Git. } }
            … fields …
            @dialog.Footer() {
                @button.Button(button.Props{Variant: button.VariantOutline, Attributes: dialog.Close(ctx)}) { Cancel }
                @button.Button(button.Props{Type: button.TypeSubmit}) { Open pull request }
            }
        </form>
    }
}
```

- Server-opened dialogs: return `dialog.Dialog(dialog.Props{DefaultOpen: true})` into an empty `#modal` target.
- One task per dialog, footer actions right-aligned: secondary then primary.
- Confirmation of destructive actions belongs in `alertdialog`.
- Preview is static (no overlay).
