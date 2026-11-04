# HoverCard

Read-only detail on hover or keyboard focus: `hovercard.HoverCard{Delay, CloseDelay}`, `hovercard.Trigger(ctx)`, `hovercard.Content{Class, Side, Align, SideOffset, AlignOffset}` from `components/hovercard` (Base UI's PreviewCard).

```templ
@hovercard.HoverCard() {
    <span class="underline decoration-dotted underline-offset-4" tabindex="0" { hovercard.Trigger(ctx)... }>{ a.Version }</span>
    @hovercard.Content(hovercard.ContentProps{Class: "w-80", Align: hovercard.AlignStart}) {
        <div class="flex flex-col gap-1 text-body-sm">…</div>
    }
}
```

- Pick by content: a short label is a `tooltip`; several lines of read-only detail (a version's image and findings, a person's card) are a hover card; anything you click inside is a `popover`. The card is a preview: what it shows also exists on the page it links to.
- Opens after 600ms and closes 300ms after the pointer leaves the trigger and the card, so the pointer can move into the card. Keyboard focus on the trigger opens it too: give a non-focusable trigger `tabindex="0"`.
- The card is portaled to `<body>` and positioned there, so a table cell's `whitespace-nowrap` and a `table-container`'s `overflow-x-auto` don't affect it. Set the width on `Content` (`w-64` by default) and break long unspaced values with `break-all` (image refs, digests).
- Load the content lazily with htmx on the trigger: `hx-get`, `hx-trigger="mouseenter once, focus once"`, `hx-target` the `Content` by its id (`HoverCard{ID}`; the id survives the portal), `hx-swap="innerHTML"`. Render a `skeleton` in `Content` until it arrives.
- Showcase: the version column on `/apps` (`VersionCard` in `ui/pages/showcase.templ`).
