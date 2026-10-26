// DatePicker glue (docs/components/DatePicker.md): copies a calendar selection
// into hidden inputs and lets htmx take it from there.
//
// Markup: a <form data-date-picker> holding the popover + calendar.Calendar and
//   <input type="hidden" data-date>                       single mode, or
//   <input type="hidden" data-date-from> + data-date-to   range mode.
// When the selection is complete the inputs get ISO dates (YYYY-MM-DD) and the
// form receives a bubbling "change", so hx-trigger="change" on the form submits it.
// The server re-renders the trigger label (e.g. via <hx-partial>); this file never
// formats dates.
(() => {
  const iso = (d) =>
    d instanceof Date
      ? `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`
      : "";

  // The popover portals its content to <body>, so the calendar is not inside the
  // form while open. The popover's trigger is: it points at the content by
  // aria-controls.
  const formOf = (el) => {
    const direct = el.closest("form[data-date-picker]");
    if (direct) return direct;
    const content = el.closest("[data-slot=popover-content]");
    const trigger = content && content.id && document.querySelector(`[aria-controls="${CSS.escape(content.id)}"]`);
    return trigger && trigger.closest("form[data-date-picker]");
  };

  document.addEventListener("calendar-select", (e) => {
    const form = e.target instanceof Element && formOf(e.target);
    if (!form) return;
    const sel = e.detail && e.detail.selected;
    const single = form.querySelector("input[data-date]");
    const from = form.querySelector("input[data-date-from]");
    const to = form.querySelector("input[data-date-to]");
    let complete = false;
    if (single) {
      single.value = iso(sel);
      complete = single.value !== "";
    } else if (from && to) {
      from.value = iso(sel && sel.from);
      to.value = iso(sel && sel.to);
      complete = from.value !== "" && to.value !== "";
    }
    if (complete) form.dispatchEvent(new Event("change", { bubbles: true }));
  });
})();
