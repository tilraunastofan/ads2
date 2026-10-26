# Calendar

A month grid for picking dates: `calendar.Calendar` from `components/calendar` (shadcn's Calendar on a port of react-day-picker).

**Consumer provides:** `Mode` (`ModeSingle`, `ModeRange`, `ModeMultiple`), `Selected` (`time.Time`, or `calendar.DateRange` for ranges), `WeekStartsOn`, `Locale`, optional `Disabled` matchers and `StartMonth` / `EndMonth`.

```templ
@calendar.Calendar(calendar.Props{
    ID:           "freeze-date",
    Mode:         calendar.ModeSingle,
    Selected:     freeze.Date,
    WeekStartsOn: calendar.Monday,
    Locale:       "is",   // Icelandic month and weekday names; "en-GB" on English pages
    Class:        "rounded-lg border",
})
```

- Weeks start on Monday (`calendar.Monday`) in every Arion app.
- `Locale` follows the page language: `is` on Icelandic pages, `en-GB` on English ones (day–month order, 24h).
- Selection is client-side: listen for the `calendar-select` event and copy `e.detail.selected` into a hidden input, then let htmx submit the form (see **DatePicker**).
- Selected day uses `primary`; today uses `muted`. Inline calendars sit in a `card` or with `border`.
