# htmx patterns

How shadcn-templ components and htmx 4 fit together in this system. htmx attributes always go through a component's `Attributes` prop (`templ.Attributes{"hx-post": url}`); the component spreads them onto its root element.

## htmx 4 defaults to remember

- **Inheritance is explicit.** Write `hx-target:inherited="#main"` on a parent to share attributes; plain attributes no longer cascade.
- **4xx/5xx responses swap.** Return 422 with the re-rendered form for validation errors. Use `hx-status:422="target:#form"` or `hx-status:5xx="target:#page-alerts"` (an empty `aria-live` region at the top of `Shell`'s `<main>`) to route errors.
- **Unhandled failures become toasts.** `request-errors.js` (loaded by `Base`) raises an error toast when the server can't be reached, and when a 5xx arrives that the element doesn't route with `hx-status` and the server didn't wrap in `<hx-partial>`; it cancels that swap so raw error text never replaces content. 4xx is left alone. On the server, wrap the handler in a recovery middleware that answers htmx requests with a toast partial and page loads with an error page (the showcase's `recoverer`).
- **One response, many targets.** Wrap extra fragments in `<hx-partial hx-target="#id">` (or `<hx-partial id="toasts">`): a row update plus a toast plus a counter, in one round trip.
- **Back/forward re-requests the page.** Make every `hx-push-url` URL render a full page when requested without `HX-Request`.
- Event names are phase-based (`htmx:before:request`, `htmx:after:swap`). Run `npx htmx.org@4.0.1 upgrade-check` on old code.

## Layout

Use `patterns.Shell` or `patterns.Base` (see **Setup**). Between them they load htmx, the component scripts and the small helpers (`toast-guard.js`, `filter-form.js`, `request-errors.js`, `dialog-form.js`, `datepicker.js`, `focus-after-swap.js`, `copy.js`, `submit-once.js`) once in `<head>`, boost the body, mount `<div id="toasts" hidden>` for toast partials, and mount the toaster. Don't hand-write the document: the pieces (order, `defer`, the `#toasts` target) are easy to get subtly wrong.

**Focus after a swap.** When a swap removes the focused element, `focus-after-swap.js` puts focus back: on the new page's `<h1>` after a boosted navigation (`PageHeader` renders it with `tabindex="-1"`), and on the first `aria-invalid="true"` field of a form re-rendered with errors. It leaves focus alone when something else placed it, such as a closing dialog. Give re-rendered forms an `id` and their invalid inputs `aria-invalid="true"` (a string).

shadcn-templ's scripts re-initialise components inside swapped content, so dialogs, menus and selects work in fragments (`parity/htmx` in the shadcn-templ repo tests exactly this).

## Extensions the shell loads

`Base` loads three htmx 4 extensions, vendored from the same release as `htmx.min.js` under `/assets/js/ext/`, and configures them in `<meta name="htmx-config" content='{"preload":{"autoBoost":true,"boostEvent":"mouseover"}}'>`. The meta comes before the htmx script (htmx reads it at start) and needs no inline script, so the CSP is unaffected. Don't add another `htmx-config` meta or load these scripts again.

| Extension | What it does |
| --- | --- |
| `hx-preload` | Fetches a boosted link's page when the pointer is over it (`mouseover`; touch starts on `touchstart`) and reuses the response if the click comes within 5 s (the default `boostTimeout`). |
| `hx-history-cache` | Keeps the last 10 pages (the default) in `sessionStorage`, so Back and Forward restore a page without a request. |
| `hx-ptag` | Per-element polling tags: an element with `hx-ptag` sends the last `HX-PTag` response header back, and the server can answer 304 when nothing changed. |

**A GET behind a boosted link may be fetched on hover.** Such a handler must have no side effects beyond audit logging, and the server should treat a preloaded request like any other. Don't put a state change behind a plain link; use a form with `hx-post`.

**Opting out.** `hx-boost="false"` on a link stops preloading it. `hx-history="false"` anywhere on a page keeps that page out of the history cache. Sign-out should clear the cache: see **Clear history on sign-out**.

**Popups on Back.** The history cache saves the page as the user leaves it, so a link clicked inside a hover card or tooltip trigger would save the card open. `history-close-popups.js` (loaded by `Base`) closes every portaled popup (not the toaster) in a cached page before it is restored. Nothing to do in apps.

**Clear history on sign-out.** `<form method="post" action="/sign-out" data-clear-history>`. `clear-history.js` (loaded by `Base`) removes the `htmx-history-*` entries from `sessionStorage` when the form is submitted, so Back after the next sign-in can't show the previous user's pages. Other `sessionStorage` keys stay. Put the attribute on every sign-out form; pages without one are unaffected.

## Recipes

**Inline validation (field):** post the form, return 422 with the same `field.Field`s re-rendered with `data-invalid` and `field.Error`.

```templ
<form hx-post="/apps" hx-target="this" hx-swap="outerHTML" hx-status:422="target:this">
    @field.Group() { … }
    @button.Button(button.Props{Type: button.TypeSubmit}) { Create application }
</form>
```

**Filter form (date range, combobox, selects → results):** `<form data-filter-form hx-get="/report/body" hx-trigger="filter-change" hx-target="#report-body" hx-swap="outerHTML">`. `filter-form.js` fires `filter-change` only when the form's named values change. Don't use `hx-trigger="change"` here: a combobox's visible text input has no name but still fires `change` when it loses focus, so the click that opens a sheet or menu in the results refetches them and the swap replaces what was just opened. `hx-trigger="change changed"` doesn't fix it either, because `changed` compares the triggering element's own value and the date picker fires `change` on the form. Lazy fragments inside the results (a sheet's body) must carry the same filter values in their URL, or they are built for the default range (the showcase's `releaseURL`). `filter-form.js` also records the form's values whenever any element inside it issues a request (`htmx:before:request`), so a live-search input inside the form (see **Live search**) doesn't make the later blur look like a changed filter: no second identical request and no duplicate `hx-push-url` history entry.

**Copy button:** `@button.Button(button.Props{Variant: button.VariantOutline, Attributes: templ.Attributes{"data-copy-target": "repo-url"}}) { @icon.Copy() <span data-label>Copy</span> }` next to `<code id="repo-url">…</code>`. `copy.js` copies the target's text and shows "Copied" / "Afritað" in `[data-label]` for two seconds. An icon-only button needs an `aria-label` that says what it copies ("Copy the Git URL").

**Submit once (sign-in, anything a second submit breaks):** `<form method="post" data-disable-on-submit data-busy-label="Signing in…">` with `<span data-label>Sign in</span>` inside its submit button. `submit-once.js` disables the submit buttons with `aria-busy` as soon as the form is submitted, and re-enables them when a boosted request finishes or the browser restores the page from the back/forward cache. Use it on a sign-in that starts an OAuth flow: a second attempt overwrites the first one's state in the session, and whichever reply comes back first then fails.

**Live search (inputgroup + table):** `hx-trigger="input changed delay:300ms, search"`, `hx-target="#rows"`, `hx-push-url="true"`. The handler renders only `table.Body` when `HX-Request` is set. Inside a `data-filter-form` form this is fine: `filter-form.js` keeps its snapshot in step with the search request, so the `change` the browser fires when the input loses focus doesn't refetch.

**Row action with confirmation (dropdownmenu → alertdialog):** the row's ⋯ menu ends with a destructive **Delete…** item that opens the row's `alertdialog` (see **DropdownMenu**). The `alertdialog.Action` carries `hx-delete`, `hx-target="#rows"`, `hx-swap="outerHTML"` and `hx-include="#search"`. The handler returns the re-rendered `table.Body` for the current query (so the last deletion shows the `empty` state) plus an `<hx-partial hx-target="#toasts">` holding a `toast.Toast`. Don't swap just the row (`hx-swap="delete"`): when the last matching row goes it leaves a bare table header instead of the `empty` state.

**Toggle a setting (switch):** `hx-post` + `hx-trigger="change from:#<ID>"` + `hx-include="#<ID>"` + `hx-swap="none"` (the switch's `change` fires on its hidden checkbox, not on the root that holds `Attributes`; see **Switch**). The response is only a `toast.Toast` partial. On error (5xx) re-render the switch in its old state.

**Polling status (card, progress, badge):** give the fragment an `ID`, `hx-get` itself with `hx-trigger="every 15s"` and `hx-swap="outerHTML"`. Stop polling by returning markup without the trigger. `every` fires while the element is hidden (an inactive tab, a collapsed accordion), and a trigger list like `intersect once, every 2s` doesn't help: each entry is its own trigger, so `every` still starts at page load. For hidden panels, render the first fragment with `hx-trigger="intersect once"`, let the server's responses carry `hx-trigger="every 2s"` while the job runs, and drop the trigger when it's done. Make the trigger a parameter of the fragment component (the showcase's `RolloutCard` does this).

**Server-opened dialog:** a button `hx-get="/apps/42/scale"` targets an empty `<div id="modal">`. The response is `dialog.Dialog(dialog.Props{DefaultOpen: true})`. On success, the form's response closes it by swapping the `#modal` contents to empty.

**Lazy panels (tabs, skeleton):** a tab's content holds `skeleton.Skeleton` plus `hx-get` with `hx-trigger="intersect once"` and `hx-swap="outerHTML"`.

**Loading state (button, spinner):** `hx-disable="this"` on the submit button and `hx-indicator="#spin"`. Show the spinner with `.htmx-request`.

## Don'ts

- No client-side state libraries next to htmx. Use a small `script` template only for purely visual behaviour shadcn-templ doesn't cover.
- Don't use `hx-confirm` for destructive actions: use `alertdialog`.
- Don't return JSON to the UI. Handlers return templ fragments, and JSON endpoints stay a separate API.
