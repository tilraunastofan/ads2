# htmx patterns

How shadcn-templ components and htmx 4 fit together in this system. htmx attributes always go through a component's `Attributes` prop (`templ.Attributes{"hx-post": url}`); the component spreads them onto its root element.

## htmx 4 defaults to remember

- **Inheritance is explicit.** Write `hx-target:inherited="#main"` on a parent to share attributes; plain attributes no longer cascade.
- **4xx/5xx responses swap.** Return 422 with the re-rendered form for validation errors. Use `hx-status:422="target:#form"` or `hx-status:5xx="target:#page-alerts"` to route errors.
- **One response, many targets.** Wrap extra fragments in `<hx-partial hx-target="#id">` (or `<hx-partial id="toasts">`): a row update plus a toast plus a counter, in one round trip.
- **Back/forward re-requests the page.** Make every `hx-push-url` URL render a full page when requested without `HX-Request`.
- Event names are phase-based (`htmx:before:request`, `htmx:after:swap`). Run `npx htmx.org@4.0.0 upgrade-check` on old code.

## Layout

Use `patterns.Shell` or `patterns.Base` (see **Setup**). Between them they load htmx and the component scripts once in `<head>`, boost the body, mount `<div id="toasts" hidden>` for toast partials, and mount the toaster. Don't hand-write the document: the pieces (order, `defer`, the `#toasts` target) are easy to get subtly wrong.

shadcn-templ's scripts re-initialise components inside swapped content, so dialogs, menus and selects work in fragments (`parity/htmx` in the shadcn-templ repo tests exactly this).

## Recipes

**Inline validation (field):** post the form, return 422 with the same `field.Field`s re-rendered with `data-invalid` and `field.Error`.

```templ
<form hx-post="/apps" hx-target="this" hx-swap="outerHTML" hx-status:422="target:this">
    @field.Group() { … }
    @button.Button(button.Props{Type: button.TypeSubmit}) { Create application }
</form>
```

**Live search (inputgroup + table):** `hx-trigger="input changed delay:300ms, search"`, `hx-target="#rows"`, `hx-push-url="true"`. The handler renders only `table.Body` when `HX-Request` is set.

**Row action with confirmation (dropdownmenu → alertdialog):** the `alertdialog.Action` carries `hx-delete`, `hx-target="#rows"`, `hx-swap="outerHTML"` and `hx-include="#search"`. The handler returns the re-rendered `table.Body` for the current query (so the last deletion shows the `empty` state) plus an `<hx-partial hx-target="#toasts">` holding a `toast.Toast`. Don't swap just the row (`hx-swap="delete"`): when the last matching row goes it leaves a bare table header instead of the `empty` state.

**Toggle a setting (switch):** `hx-post` + `hx-trigger="change"` + `hx-swap="none"`. The response is only a `toast.Toast` partial. On error (5xx) re-render the switch in its old state.

**Polling status (card, progress, badge):** give the fragment an `ID`, `hx-get` itself with `hx-trigger="every 15s"` and `hx-swap="outerHTML"`. Stop polling by returning markup without the trigger. `every` fires while the element is hidden (an inactive tab, a collapsed accordion), and a trigger list like `intersect once, every 2s` doesn't help: each entry is its own trigger, so `every` still starts at page load. For hidden panels, render the first fragment with `hx-trigger="intersect once"`, let the server's responses carry `hx-trigger="every 2s"` while the job runs, and drop the trigger when it's done. Make the trigger a parameter of the fragment component (the showcase's `RolloutCard` does this).

**Server-opened dialog:** a button `hx-get="/apps/42/scale"` targets an empty `<div id="modal">`. The response is `dialog.Dialog(dialog.Props{DefaultOpen: true})`. On success, the form's response closes it by swapping the `#modal` contents to empty.

**Lazy panels (tabs, skeleton):** a tab's content holds `skeleton.Skeleton` plus `hx-get` with `hx-trigger="intersect once"` and `hx-swap="outerHTML"`.

**Loading state (button, spinner):** `hx-disable="this"` on the submit button and `hx-indicator="#spin"`. Show the spinner with `.htmx-request`.

## Don'ts

- No client-side state libraries next to htmx. Use a small `script` template only for purely visual behaviour shadcn-templ doesn't cover.
- Don't use `hx-confirm` for destructive actions: use `alertdialog`.
- Don't return JSON to the UI. Handlers return templ fragments, and JSON endpoints stay a separate API.
