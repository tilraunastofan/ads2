# StatusBadge

A state as a pill with icon and word, mapped the same way in every app. `patterns.StatusBadge(Status, label)` from `github.com/tilraunastofan/ads2/patterns` (copy the package into your app like `components/`).

**Consumer provides:** one `patterns.Status` — `StatusPositive` (healthy, synced, succeeded), `StatusWarning` (progressing, needs attention), `StatusDestructive` (failed, degraded, down, out of sync), `StatusNeutral` (unknown, pending, archived) — and the label in the page language.

```templ
func healthStatus(h string) patterns.Status {
    switch h {
    case "Healthy":
        return patterns.StatusPositive
    case "Progressing":
        return patterns.StatusWarning
    case "Degraded", "Missing", "Failed":
        return patterns.StatusDestructive
    }
    return patterns.StatusNeutral
}

@patterns.StatusBadge(healthStatus(a.Health), a.Health)
```

- Map your domain states onto `Status` in one helper per app; never pick badge variants per page.
- Every status has its own icon (check, triangle, alert circle, dashed circle), so colour never carries the state alone.
- For metadata that is not a state (versions, environments) use a plain `badge`.
