// Clear history: a form marked data-clear-history (sign-out) removes the pages
// hx-history-cache keeps in sessionStorage when it is submitted, so the next user of
// this browser tab cannot press Back and see the previous user's pages.
//
// hx-history-cache writes "htmx-history-index" and one "htmx-history-<id>" per page;
// the prefix below covers both. Other sessionStorage keys are left alone. A plain
// file in /assets/js because the CSP (script-src 'self') blocks inline handlers.
(() => {
  const PREFIX = "htmx-history-";

  document.addEventListener("submit", (e) => {
    if (!(e.target instanceof Element) || !e.target.closest("form[data-clear-history]")) return;
    try {
      const keys = [];
      for (let i = 0; i < sessionStorage.length; i++) keys.push(sessionStorage.key(i));
      for (const key of keys) if (key.startsWith(PREFIX)) sessionStorage.removeItem(key);
    } catch {
      // sessionStorage is unavailable: the cache could not have been written either.
    }
  });
})();
