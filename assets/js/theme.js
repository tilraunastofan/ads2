// Dark mode is opt-in by class on <html>. The first visit follows the OS setting,
// live; after a click on the toggle its choice is remembered per browser and wins.
// Load this before first paint (a plain <script> in <head>) so there is no flash.
// Toggles ([data-theme-toggle]) get aria-pressed="true" while dark is on.
(() => {
  const root = document.documentElement;
  const os = matchMedia("(prefers-color-scheme: dark)");
  const saved = () => {
    try { return localStorage.getItem("theme"); } catch (_) { return null; }
  };
  const sync = () => {
    const dark = root.classList.contains("dark");
    document.querySelectorAll("[data-theme-toggle]").forEach((b) => b.setAttribute("aria-pressed", String(dark)));
  };
  const apply = (dark) => {
    root.classList.toggle("dark", dark);
    sync();
  };

  const s = saved();
  apply(s ? s === "dark" : os.matches);

  os.addEventListener("change", (e) => { if (!saved()) apply(e.matches); });
  document.addEventListener("DOMContentLoaded", sync);
  // Boosted navigation swaps the body, which re-renders the toggle.
  document.addEventListener("htmx:after:swap", sync);
  document.addEventListener("click", (e) => {
    if (!(e.target instanceof Element) || !e.target.closest("[data-theme-toggle]")) return;
    const dark = !root.classList.contains("dark");
    apply(dark);
    try { localStorage.setItem("theme", dark ? "dark" : "light"); } catch (_) {}
  });
})();
