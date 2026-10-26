# Sheet

A side panel over the page for detail and secondary tasks: `sheet.Sheet`, `sheet.Trigger(ctx)`, `Content{Side}`, `Header`, `Title`, `Description`, `Footer`, `sheet.Close(ctx)` from `components/sheet`.

**Consumer provides:** the trigger (any button, through `Attributes: sheet.Trigger(ctx)`), `Side` (`SideRight` by default), title, description and body.

```templ
@sheet.Sheet(sheet.Props{ID: "details-" + a.Name}) {
    @button.Button(button.Props{Variant: button.VariantGhost, Size: button.SizeSm, Attributes: sheet.Trigger(ctx)}) { Details }
    @sheet.Content() {
        @sheet.Header() {
            @sheet.Title() { { a.Name } }
            @sheet.Description() { Release { a.Version } · { a.Env } }
        }
        <div hx-get={ "/apps/" + a.Name + "/details" } hx-trigger="intersect once" hx-swap="outerHTML" class="px-4">
            @skeleton.Skeleton(skeleton.Props{Class: "h-32"})
        </div>
        @sheet.Footer() {
            @button.Button(button.Props{Variant: button.VariantOutline, Attributes: sheet.Close(ctx)}) { Close }
        }
    }
}
```

- Row details, filters and logs: anything people look at beside the list without losing their place. A task that must be completed or cancelled is a `dialog`.
- Load the body lazily with `hx-trigger="intersect once"` (the panel is hidden until opened).
- Right side for detail, left only for navigation on small screens (the sidebar does this itself).
- The preview shows the open panel inside a frame, without the overlay.
