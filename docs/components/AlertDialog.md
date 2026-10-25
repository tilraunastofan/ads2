# AlertDialog

A confirmation that must be answered: `alertdialog.AlertDialog`, `alertdialog.Trigger(ctx)`, `Content{Size}`, `Header`, `Media`, `Title`, `Description`, `Footer`, `Cancel`, `Action{Variant}` from `components/alertdialog`.

```templ
@alertdialog.AlertDialog() {
    @button.Button(button.Props{Variant: button.VariantDestructive, Attributes: alertdialog.Trigger(ctx)}) { Delete }
    @alertdialog.Content() {
        @alertdialog.Header() {
            @alertdialog.Title() { Delete namespace “{ ns }”? }
            @alertdialog.Description() { All workloads, secrets and volumes are removed. This cannot be undone. }
        }
        @alertdialog.Footer() {
            @alertdialog.Cancel() { Cancel }
            @alertdialog.Action(alertdialog.ActionProps{Variant: button.VariantDestructive,
                Attributes: templ.Attributes{"hx-delete": "/namespaces/" + ns, "hx-target": "#ns-rows", "hx-swap": "outerHTML", "hx-include": "#ns-search"}}) { Delete namespace }
        }
    }
}
```

- Use instead of htmx `hx-confirm` for anything irreversible: it is styled, accessible, and names the object.
- The action button repeats the verb and object ("Delete namespace"), never "OK".
- The delete handler returns the refreshed rows (empty state included) plus a toast partial; see **htmx patterns**.
