// Copy buttons: a button with data-copy-target="<id>" copies the text of the
// element with that id, then shows "Copied" (or "Copy failed") for two seconds in
// its [data-label] element, or in the button itself. Words follow <html lang>:
// "Afritað" / "Afritun mistókst" on Icelandic pages. Loaded once by patterns.Base;
// a plain file because the CSP (script-src 'self') blocks inline handlers.
(() => {
  const words = () =>
    document.documentElement.lang.startsWith("is")
      ? { ok: "Afritað", failed: "Afritun mistókst" }
      : { ok: "Copied", failed: "Copy failed" };

  document.addEventListener("click", async (e) => {
    const btn = e.target.closest("[data-copy-target]");
    if (!btn) return;
    const src = document.getElementById(btn.dataset.copyTarget);
    if (!src) return;
    const label = btn.querySelector("[data-label]") || btn;
    label.dataset.idle ??= label.textContent;
    try {
      await navigator.clipboard.writeText(src.textContent);
      label.textContent = words().ok;
    } catch {
      label.textContent = words().failed;
    }
    setTimeout(() => (label.textContent = label.dataset.idle), 2000);
  });
})();
