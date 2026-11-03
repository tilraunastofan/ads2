// Forms in a dialog or sheet: <form data-close-on-success hx-post=…> closes its
// dialog once the request succeeds (2xx), so a failed request (4xx re-render, 5xx
// toast) leaves it open with the user's input.
//
// Closing on the submit button's click doesn't work: the dialog closes before the
// form submits, and the request never goes out. hx-on would need eval, which the
// CSP blocks. So this waits for htmx's answer and clicks the dialog's own close
// control, which keeps focus return and the close animation.
document.addEventListener("htmx:after:request", (e) => {
  const form = e.target instanceof Element && e.target.closest("form[data-close-on-success]");
  const status = e.detail?.ctx?.response?.status;
  if (!form || !status || status < 200 || status >= 300) return;
  const popup = form.closest('[role="dialog"], [role="alertdialog"]');
  popup?.querySelector("[data-templ-dialog-close]")?.click();
});
