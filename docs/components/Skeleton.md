# Skeleton

Placeholder while a region loads: `skeleton.Skeleton` from `components/skeleton` (`bg-muted`, pulsing).

```templ
<div hx-get="/apps/payments-api/metrics" hx-trigger="load" hx-swap="outerHTML">
    @skeleton.Skeleton(skeleton.Props{Class: "h-48 w-full rounded-xl"})
</div>
```

- Match the shape of what loads (card, row, avatar). Use for lazy `hx-trigger="load"` / `intersect once` regions.
