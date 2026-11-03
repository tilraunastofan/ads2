// Feedback for htmx requests that fail without a server-rendered answer.
//
// - The server can't be reached (offline, timeout, connection refused): htmx fires
//   htmx:error and swaps nothing, so the click would look ignored. Raise a toast.
// - A 5xx arrives that neither the element routes (hx-status:5xx / hx-status:503 …)
//   nor the server wrapped in an <hx-partial> (the showcase's recovery middleware
//   does): htmx 4 would swap the raw body ("Internal Server Error", a proxy's 502
//   page) into the target. Cancel that swap and raise a toast instead.
//
// 4xx stays untouched: 422 re-renders a form, 404 for a fragment is the handler's
// call. Messages follow the page language (<html lang>).
(() => {
  const words = {
    en: {
      offlineTitle: "Request failed",
      offline: "The server could not be reached. Check your connection and try again.",
      serverTitle: "Something went wrong",
      server: "The server could not complete the request. Try again; if it keeps failing, report it with the time.",
    },
    is: {
      offlineTitle: "Beiðni mistókst",
      offline: "Ekki náðist samband við þjóninn. Athugaðu tenginguna og reyndu aftur.",
      serverTitle: "Eitthvað fór úrskeiðis",
      server: "Þjónninn gat ekki lokið beiðninni. Reyndu aftur; ef það gengur ekki, tilkynntu það með tímanum.",
    },
  };
  const t = () => words[(document.documentElement.lang || "en").slice(0, 2)] || words.en;
  const toast = (title, description) => window.templ?.toast?.add({ title, description, type: "error" });

  // An hx-status rule on the source element or an ancestor (hx-status:5xx:inherited).
  const routed = (el, status) => {
    for (let node = el; node instanceof Element; node = node.parentElement) {
      for (const a of node.attributes) {
        if (!a.name.startsWith("hx-status:")) continue;
        const code = a.name.slice("hx-status:".length).split(":")[0];
        const own = node === el || a.name.endsWith(":inherited");
        if (own && (code === String(status) || code === `${String(status)[0]}xx`)) return true;
      }
    }
    return false;
  };

  document.addEventListener("htmx:error", (e) => {
    if (!e.detail?.ctx) return; // an exception in a handler, not a failed request
    const w = t();
    toast(w.offlineTitle, w.offline);
  });

  document.addEventListener("htmx:before:swap", (e) => {
    const ctx = e.detail?.ctx;
    const status = ctx?.response?.status;
    if (!status || status < 500) return;
    if (routed(ctx.sourceElement, status) || /<hx-partial[\s>]/i.test(ctx.text || "")) return;
    e.preventDefault();
    const w = t();
    toast(w.serverTitle, w.server);
  });
})();
