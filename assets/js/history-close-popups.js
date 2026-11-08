// Close popups in restored pages: hx-history-cache saves a copy of the page
// as the user leaves it. Clicking a link inside a hover card or tooltip trigger
// saves the page with that card open, and Back would restore it open with no
// script state behind it, so it stays until hovered again. Before a cached page
// is restored, every portaled popup in the copy (hover card, tooltip, popover,
// menu, select, dialog) is put back in its closed, server-rendered state.
// The toast portal is always visible and is left alone. A plain file in
// /assets/js because the CSP (script-src 'self') blocks inline handlers.
(() => {
  function closePopups(html) {
    const doc = new DOMParser().parseFromString(html, "text/html");
    const root = doc.body;
    let changed = false;
    for (const portal of root.querySelectorAll("[data-base-ui-portal]")) {
      if (portal.getAttribute("data-slot") === "toast-portal") continue;
      if (portal.querySelector("[data-templ-open]")) continue; // controlled open
      for (const el of [portal, ...portal.querySelectorAll("[data-open]")]) {
        if (el.hasAttribute("data-open")) {
          el.removeAttribute("data-open");
          el.setAttribute("data-closed", "");
          changed = true;
        }
      }
      for (const el of [portal, ...portal.children]) {
        if (!el.hidden) {
          el.hidden = true;
          changed = true;
        }
      }
    }
    for (const trigger of root.querySelectorAll("[data-popup-open]")) {
      trigger.removeAttribute("data-popup-open");
      changed = true;
    }
    return changed ? root.outerHTML : html;
  }

  document.addEventListener("htmx:history:cache:hit", (e) => {
    const item = e.detail && e.detail.item;
    if (!item || typeof item.content !== "string") return;
    try {
      item.content = closePopups(item.content);
    } catch {
      // Leave the page as saved rather than break Back.
    }
  });
})();
