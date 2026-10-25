# Switch

An on/off setting that applies immediately: `switch.Switch` from `components/switch` (import as `switchcomp` or similar — `switch` is a Go keyword).

**Consumer provides:** `ID`, `Name`, `DefaultChecked` (or `Checked *bool`), `Size` (`SizeDefault`, `SizeSm`), a label.

```templ
@switchcomp.Switch(switchcomp.Props{
    ID: "autosync", Name: "autosync", DefaultChecked: app.AutoSync,
    Attributes: templ.Attributes{"hx-post": "/apps/" + app.Name + "/autosync", "hx-trigger": "change", "hx-swap": "none"},
})
```

- Saves on change (htmx `hx-trigger="change"`) and confirms with a `toast`; never inside a form with a Save button.
- Unchecked track uses `input`, checked uses `primary`.
