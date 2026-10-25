# Progress

Determinate progress: `progress.Progress{Value, Max}` with `Label`, `Value`, `Track`, `Indicator` from `components/progress`.

```templ
// trigger: "intersect once" on the page, "every 2s" in poll responses.
// The htmx attributes are dropped once ready == desired, which stops polling.
templ RolloutCard(ready, desired int, trigger string) {
    <div id="rollout" if ready < desired { hx-get="/rollout" hx-trigger={ trigger } hx-swap="outerHTML" }>
        @progress.Progress(progress.Props{Value: ready, Max: desired})
    </div>
}
```

- Inside a tab or accordion, start with `intersect once` and switch to `every 2s` in the responses. `intersect once, every 2s` does not wait: `every` starts at page load.
- Rollouts, imports, batch jobs. Stop polling when complete by returning the final markup without `hx-trigger` (or respond 286).
- Unknown duration → `spinner`.
