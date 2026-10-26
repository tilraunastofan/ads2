# DatePicker

A button that opens a calendar in a popover. Not a separate component: compose `popover.Popover`, `popover.Trigger(ctx)`, `popover.Content` and `calendar.Calendar` (`shadcn-templ add popover calendar`).

**Consumer provides:** the trigger label (the formatted date or range, or a placeholder), `Mode`, hidden inputs that carry the value to the server, and the small script that fills them.

```templ
<form id="range" hx-get="/metrics" hx-target="#metrics" hx-select="#metrics" hx-trigger="change" hx-push-url="true">
    <input type="hidden" name="from" value={ r.From.Format(time.DateOnly) }/>
    <input type="hidden" name="to" value={ r.To.Format(time.DateOnly) }/>
    @popover.Popover(popover.Props{ID: "range-picker"}) {
        @button.Button(button.Props{Variant: button.VariantOutline, Class: "w-[260px] justify-between font-normal", Attributes: popover.Trigger(ctx)}) {
            <span class="flex items-center gap-2">@icon.Calendar() { patterns.DateRange(r.From, r.To, "en") }</span>
            @icon.ChevronDown()
        }
        @popover.Content(popover.ContentProps{Class: "w-auto p-0", Align: popover.AlignStart}) {
            @calendar.Calendar(calendar.Props{Mode: calendar.ModeRange, Selected: calendar.DateRange{From: r.From, To: r.To}, WeekStartsOn: calendar.Monday, NumberOfMonths: 2})
        }
    }
</form>
<script src="/assets/js/datepicker.js" defer></script>
```

`assets/js/datepicker.js` copies `calendar-select` into the hidden `from` / `to` inputs (ISO `YYYY-MM-DD`) and dispatches `change` on the form once a range is complete, so htmx re-renders the page region. The server owns the label: it re-renders the trigger with the formatted range.

- The popover portals its content to `<body>` while open, so the calendar is **not** a DOM descendant of your form. `datepicker.js` therefore finds the form through the popover trigger (`aria-controls` points at the content id), which is inside the form. Keep the `popover.Trigger(ctx)` attributes on a button in the same `form[data-date-picker]`, and give the `Popover` an `ID`.
- The popover stays open after a selection, so a range can be adjusted; the first click inside an existing range moves its nearest end. Escape or an outside click closes it.
- Labels come from `patterns.Date` / `patterns.DateRange`: `20–24 Oct 2026` on English pages, `20.–24. október 2026` on Icelandic ones. Values on the wire are ISO dates (`YYYY-MM-DD`).
- Ranges for reports default to the last full week or month, never an empty picker.
- The preview shows the popover open below its trigger.
