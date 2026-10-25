// Dark mode is opt-in by class on <html>. The first visit follows the OS setting;
// after that the toggle's choice is remembered per browser. Load this before
// first paint (a plain <script> in <head>) so there is no flash.
try {
  const saved = localStorage.getItem("theme");
  const dark = saved ? saved === "dark" : matchMedia("(prefers-color-scheme: dark)").matches;
  if (dark) document.documentElement.classList.add("dark");
} catch (_) {}
document.addEventListener("click", (e) => {
  if (!e.target.closest("[data-theme-toggle]")) return;
  const dark = document.documentElement.classList.toggle("dark");
  try { localStorage.setItem("theme", dark ? "dark" : "light"); } catch (_) {}
});
