# Toast freezes the tab on page change

**Status:** worked around in this repo (`assets/js/toast-guard.js`, commit `f7e017e`). Not fixed upstream. Reported 3 Nov 2026 as [axadrn/shadcn-templ#656](https://github.com/axadrn/shadcn-templ/issues/656), with the [standalone reproduction](https://github.com/jakobjs/shadcn-templ-toast-freeze-repro) and the fix under **Fixing it upstream**.
**Affects:** shadcn-templ's toast (`components/toast/toast.js`) in every v2 release: reproduced on `v2.0.0-beta.15`; the file is byte-identical in `beta.16`, `v2.0.0` and `main` as of 3 Nov 2026, so upgrading does not fix it. ads2 runs `v2.0.0` (since 3 Nov 2026); after that upgrade the three checks under **After every shadcn-templ update** passed in Chrome and the guard check printed `true`. Independent of the htmx version (4.0.0 and 4.0.1 both freeze). shadcn's own React toast is not affected (see **Upstream analysis**).
**Severity:** high. The tab hangs until it is closed, and it takes only an ordinary click.

## What users see

A toast is showing ("Restarting pods", "Application created" …). Before it disappears, the user clicks a sidebar link or presses Back. The URL changes, the page doesn't, and the tab stops responding. Closing the tab is the only way out.

## Reproduce

Standalone: https://github.com/jakobjs/shadcn-templ-toast-freeze-repro is a fresh `shadcn-templ init` (v2.0.0) with only the toast, a boosted body and two pages; `go run .` and follow its README. It freezes on a clean project, so the cause is in shadcn-templ, not in this repo.

In this repo:

Without the workaround (remove the `toast-guard.js` `<script>` from `patterns.Base`, or test an app that never had it):

1. Open `/apps`.
2. Open a row's menu (⋯) and choose **Restart**. A toast appears.
3. While the toast is visible, click **Settings** in the sidebar.

The tab freezes. It doesn't matter which action raised the toast or which link is clicked. Any boosted navigation or history restore that swaps `<body>` triggers it. It happens on htmx 4.0.0 and 4.0.1, so it is not an htmx regression.

## Cause

`hx-boost` swaps the whole `<body>`, which removes the toast viewport. shadcn-templ's lifecycle then calls the viewport's `destroy()` (`components/toast/toast.js` around line 1116). It unsubscribes from the toast manager and clears the timers, but it does **not** disconnect the two observers each toast creates when it mounts (lines 716–717):

```js
r.observers = [
  new ResizeObserver(() => r.recalculateHeight()),
  new MutationObserver(() => r.recalculateHeight()),
];
```

The observers stay attached to the now-detached toast. `recalculateHeight()` updates the toast store, the store schedules a re-render, and the re-render writes to the toast's DOM, which fires the `MutationObserver` again. That loop runs in microtasks, so the browser never gets to paint or handle input again.

The component does disconnect a toast's observers when that toast is closed normally (line 1034). Only the path where the whole viewport goes away misses it.

**Evidence:** with `MutationObserver` wrapped in a counter, a boosted navigation with one toast showing produced 2.000+ callbacks within seconds (zero before the click). With the observers disconnected first, the count stayed at 0 and the navigation completed.

## Upstream analysis

Compared on 3 Nov 2026, from source.

| Code | Observers cleaned up when the toaster goes away? | Affected |
| --- | --- | --- |
| shadcn `apps/v4/registry/bases/base/ui/toast.tsx` (the Base UI variant) | Adds no observers or effects of its own; composes Base UI's `Toast` primitives. | No |
| Base UI `packages/react/src/toast/content/ToastContent.tsx` (1.6.0, which shadcn-templ ports, and `master`) | Yes. The observers live in a layout effect whose cleanup disconnects both, and React runs it whenever a toast unmounts, including when its viewport or provider unmounts. | No |
| shadcn-templ `components/toast/toast.js` (`beta.15`, `beta.16`, `v2.0.0`, `main`; SHA-256 `300c0829…` for all four) | No. Per-toast cleanup runs only when one toast leaves the list (`renderList`, line 1034); the viewport's `destroy()` (line 1116) unsubscribes and clears timers but never runs it. | Yes |

So this is a gap in the port, not in the original. React tears down every child of an unmounted tree; the vanilla port has to do that by hand in `destroy()`, and doesn't. React apps also rarely unmount the toaster at all (it sits in the root layout and survives client-side navigation), while `hx-boost` replaces `<body>`, and with it the toaster, on every navigation.

`components/baseui/lifecycle.js`, which calls `destroy()`, is also identical across these versions. Upstream issue #514 ("Tolerate HTMX history restore with Toast timers", v1 line, closed) is related but different: a restarted timer, not a freeze.

## Workaround in this repo

`assets/js/toast-guard.js`, loaded by `patterns.Base`, listens for `htmx:before:swap`. When the swap target is `<body>`, it disconnects every toast's observers before htmx removes the old body:

```js
viewport._templToastProvider?.toasts?.forEach((toast) => {
  toast.observers?.forEach((o) => o.disconnect());
});
```

Toasts still disappear with the old page, as they would without the bug, and toasts on the new page work normally.

## Risk: the workaround depends on internals

`_templToastProvider`, its `toasts` map and each toast's `observers` array are private to the component, not a public API. shadcn-templ is a beta and Dependabot proposes updates weekly. If an update renames or restructures any of these:

- the guard's optional chaining (`?.`) means it **does nothing, silently**: no error, no console message;
- the freeze comes back;
- **no Go test notices.** The bug lives in browser JavaScript, and `go test` only sees server-rendered HTML.

## After every shadcn-templ update

Until the bug is fixed upstream, check this by hand in Chrome whenever `github.com/axadrn/shadcn-templ/v2` changes, or `components/toast/` is re-added:

1. `task dev`, open http://localhost:7331/apps.
2. Raise a toast (row menu → **Restart**), and while it shows click **Settings**. The Settings page must load.
3. Raise a toast on Deployment freeze (code `123456`, **Lift freeze**), and while it shows press Back. Applications must load.

For a quick look at whether the guard still reaches the internals, run this in the console while a toast is showing. It must print `true`:

```js
[...document.querySelectorAll('[data-slot="toast-viewport"]')]
  .some((v) => v._templToastProvider?.toasts?.size > 0 && [...v._templToastProvider.toasts.values()].every((t) => Array.isArray(t.observers)))
```

If it prints `false` with a toast visible, the component changed: read the new `toast.js`, update `toast-guard.js`, or remove it if `destroy()` now disconnects the observers itself.

## Fixing it upstream

Report the bug to shadcn-templ (`github.com/axadrn/shadcn-templ`, bug form, title prefix `[bug]: `). The fix mirrors what React does on unmount: in the viewport's `destroy()`, run the same per-toast cleanup `renderList` already runs for a removed toast:

```js
destroy(viewport) {
  const provider = viewport._templToastProvider;
  if (!provider) return;
  provider.toasts.forEach((r) => {            // missing today
    r.observers.forEach((o) => o.disconnect());
    r.dragAbortController?.abort();
  });
  provider.unsubscribe();
  provider.store.clearTimers();
},
```

Reported as [#656](https://github.com/axadrn/shadcn-templ/issues/656). Follow that issue; when a release closes it, do the removal steps below.

Once a release contains it: update the module, run the checks above **without** `toast-guard.js`, then delete the script, its `<script>` tag in `patterns/shell.templ`, its rows in `CLAUDE.md`, and this file.
