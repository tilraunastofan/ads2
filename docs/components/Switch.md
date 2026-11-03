# Switch

An on/off setting that applies immediately: `switch.Switch` from `components/switch` (import as `switchcomp` or similar — `switch` is a Go keyword).

**Consumer provides:** `ID`, `Name`, `DefaultChecked` (or `Checked *bool`), `Size` (`SizeDefault`, `SizeSm`), a label.

```templ
@switchcomp.Switch(switchcomp.Props{
    ID: "autosync", Name: "autosync", DefaultChecked: app.AutoSync,
    Attributes: templ.Attributes{
        "hx-post": "/apps/" + app.Name + "/autosync", "hx-swap": "none",
        "hx-trigger": "change from:#autosync", "hx-include": "#autosync",
    },
})
```

- Saves on change and confirms with a `toast`; never inside a form with a Save button.
- `Attributes` land on the switch root, but the `change` event fires on the hidden checkbox next to it, which does not bubble through the root. A plain `hx-trigger="change"` never fires: use `change from:#<ID>`, and `hx-include="#<ID>"` so the request carries the checkbox value (present only when on).
- Name it for screen readers: the label's `for` reaches the hidden checkbox, not the `role="switch"` element they announce. Give the label an `ID` and the switch `"aria-labelledby": "<label ID>"` in `Attributes`.
- Unchecked track uses `input`, checked uses `primary`.
