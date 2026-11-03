# Stat

A KPI card: metric name, the value as the largest thing in the card, and under it the change with sign, arrow and colour, followed by `Since` ("vs last week", "miðað við 3F 2025"). `patterns.Stat(StatProps)` and the grid `patterns.StatRow()` from `github.com/tilraunastofan/ads2/patterns` (imported from the ads2 module, see `docs/setup.md`).

**Consumer provides:** `Label`, `Value` (already formatted with `patterns.Int`, `Dec`, `Pct`), `Delta` (formatted with `Signed`, empty hides the badge), and two independent flags: `Up` (the number went up: drives the arrow) and `Good` (the change is good: drives the colour).

```templ
@patterns.StatRow() {
    @patterns.Stat(patterns.StatProps{Label: "Deploys this week", Value: patterns.Int(128), Delta: patterns.SignedPct(18, 0), Good: true, Up: true})
    @patterns.Stat(patterns.StatProps{Label: "Change failure rate", Value: patterns.Pct(2.3, 1), Delta: patterns.Signed(-0.8, 1) + " pp", Good: true, Up: false})
    @patterns.Stat(patterns.StatProps{Label: "Mean time to restore", Value: "41 min", Delta: "+6 min", Good: false, Up: true})
}
```

- Three or four per row, directly under the `PageHeader`.
- A falling failure rate is down and good: never derive colour from the sign.
- Values in Arion format in every language; units follow the page language ("41 min" / "41 mín.").
