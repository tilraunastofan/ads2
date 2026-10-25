# Progress

Determinate progress: `progress.Progress{Value, Max}` with `Label`, `Value`, `Track`, `Indicator` from `components/progress`.

```templ
<div id="rollout" hx-get={ rolloutURL } hx-trigger="every 2s" hx-swap="outerHTML">
    @progress.Progress(progress.Props{Value: r.Ready, Max: r.Desired})
</div>
```

- Rollouts, imports, batch jobs. Stop polling when complete by returning the final markup without `hx-trigger` (or respond 286).
- Unknown duration → `spinner`.
