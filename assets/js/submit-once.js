// Submit once: a form marked data-disable-on-submit greys out its submit button as
// soon as it is submitted, so a second click cannot start a second request. For the
// Entra sign-in that matters beyond looks: a second attempt overwrites the first one's
// state and nonce in the session, and whichever reply comes back first then fails.
//
// The button's label element ([data-label]) shows the form's data-busy-label while
// the browser is busy. A plain file in /assets/js because the CSP (script-src 'self')
// blocks inline handlers, including htmx's hx-on.
(() => {
  const busy = (form, on) => {
    for (const btn of form.querySelectorAll('button[type="submit"]')) {
      btn.disabled = on;
      btn.toggleAttribute("aria-busy", on);
      const label = btn.querySelector("[data-label]");
      if (!label) continue;
      if (on) {
        label.dataset.idle ??= label.textContent;
        label.textContent = form.dataset.busyLabel || label.dataset.idle;
      } else if (label.dataset.idle !== undefined) {
        label.textContent = label.dataset.idle;
      }
    }
  };
  const forms = () => document.querySelectorAll("form[data-disable-on-submit]");

  document.addEventListener("submit", (e) => {
    const form = e.target.closest("form[data-disable-on-submit]");
    // Disable after the browser has taken the submitter, so the submission goes ahead.
    if (form) setTimeout(() => busy(form, true));
  });
  // Back from Entra (or any page) can restore this one from the back/forward cache
  // with the button still disabled; make it usable again.
  window.addEventListener("pageshow", (e) => {
    if (e.persisted) forms().forEach((f) => busy(f, false));
  });
  // A boosted (htmx) submit that fails leaves the page in place: re-enable its form.
  document.addEventListener("htmx:after:request", (e) => {
    const form = e.target.closest?.("form[data-disable-on-submit]");
    if (form) busy(form, false);
  });
})();
