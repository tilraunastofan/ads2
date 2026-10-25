# Toast

Transient confirmation: `toast.Toaster` (once in the layout) and `toast.Toast{Title, Description, Type, Timeout}` from `components/toast`. A `toast.Toast` rendered in any htmx response is picked up by `toast.js` and shown.

```templ
// layout.templ, once
@toast.Toaster()

// handler response fragment (e.g. via <hx-partial> or appended to the swap)
@toast.Toast(toast.Props{Title: "Promoted to production", Description: app.Name + " " + v + " is rolling out.", Type: toast.TypeSuccess})
```

- Types: `TypeSuccess`, `TypeInfo`, `TypeWarning`, `TypeError`, `TypeLoading`; each draws its own icon.
- Confirms what just happened; never the only place an error is shown (use `field.Error` or `alert`).
- Preview is static; the live toast stacks bottom-right.
