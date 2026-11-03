// Focus after htmx swaps. When a swap removes the element that had focus, the
// browser drops focus to <body>: a keyboard user starts again at the top of the
// page, and a screen reader says nothing about what changed. This puts focus where
// the user can carry on:
//
// - Boosted navigation (the body is swapped): the new page's <h1>, so a screen
//   reader announces the page and Tab continues into its content. PageHeader's h1
//   carries tabindex="-1" for this; without one, #main (Shell's skip target).
// - A form re-rendered with errors (422, hx-swap="outerHTML"): its first invalid
//   field, so the error that aria-describedby ties to it is read out.
//
// It only acts when focus was really lost: on <body>, or (for a form) on a container
// of the new form, which is where a dialog's focus trap puts it. If something else
// placed focus (a dialog closing returns it to its trigger), it leaves it alone.
(() => {
  if (window.arionFocusAfterSwap) return; // loaded twice: keep one listener
  window.arionFocusAfterSwap = true;

  const lost = () => !document.activeElement || document.activeElement === document.body;

  const focusable = (el) =>
    el && !el.disabled && el.getClientRects().length > 0 && (el.tabIndex >= 0 || el.hasAttribute("tabindex"));

  const firstInvalid = (root) =>
    [...(root?.querySelectorAll('[aria-invalid="true"]') ?? [])].find(
      (el) => focusable(el) && el.matches("input, select, textarea, [contenteditable], [tabindex]"),
    );

  document.addEventListener("htmx:before:swap", (e) => {
    const ctx = e.detail?.ctx;
    const active = document.activeElement;
    if (!ctx || !active || active === document.body) return;
    const boosted = ctx.target === document.body;
    const formId = active.closest("form[id]")?.id;

    // The swap runs after this event, and its event fires on the replaced (and so
    // disconnected) element, which never reaches document. So wait a few frames for
    // the focused element to leave the page instead.
    let frames = 0;
    const settle = () => {
      if (active.isConnected && frames++ < 30) return requestAnimationFrame(settle);
      if (active.isConnected) return;
      // A dialog's focus trap reacts to the removal too and puts focus on the
      // dialog a frame later; act after it, or it takes focus straight back.
      requestAnimationFrame(() => requestAnimationFrame(place));
    };
    const place = () => {
      if (boosted) {
        if (!lost()) return;
        const target = document.querySelector("#main h1[tabindex]") || document.getElementById("main");
        target?.focus({ preventScroll: true });
        return;
      }
      const form = formId && document.getElementById(formId);
      const trapped = form && document.activeElement?.contains(form);
      if (form && (lost() || trapped)) firstInvalid(form)?.focus();
    };
    requestAnimationFrame(settle);
  });
})();
