# Textarea

Multi-line text entry: `textarea.Textarea` from `components/textarea`; grows with content (`field-sizing-content`).

**Consumer provides:** `Name`, `Placeholder`, `Rows`, `Value`, an `ID` for its label.

```templ
@textarea.Textarea(textarea.Props{ID: "reason", Name: "reason", Rows: 3, Placeholder: "Why is this change needed?"})
```

- Change requests, incident notes, comments. Show a character limit with `field.Description` when one exists.
