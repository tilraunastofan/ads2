// Workaround for a freeze in the vendored shadcn-templ toast (components/toast/toast.js,
// which we don't edit): a boosted navigation swaps <body> while a toast is showing,
// the viewport's destroy() unsubscribes and clears timers but leaves each toast's
// ResizeObserver and MutationObserver connected, and they re-render the detached
// toast in an endless microtask loop. The tab hangs until it is closed.
//
// Before htmx swaps the body, disconnect those observers. The toasts go away with the
// old body, as they would without the bug. If the component's internals change, this
// does nothing. Remove it once upstream disconnects the observers in destroy().
document.addEventListener("htmx:before:swap", (e) => {
  if (e.detail?.ctx?.target !== document.body) return;
  document.querySelectorAll('[data-slot="toast-viewport"]').forEach((viewport) => {
    viewport._templToastProvider?.toasts?.forEach((toast) => {
      toast.observers?.forEach((o) => o.disconnect());
    });
  });
});
