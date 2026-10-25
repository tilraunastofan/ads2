# Badge

Short status or metadata label: `badge.Badge` from `components/badge` (`<span>`, or `<a>` with `Href`).

**Consumer provides:** one to three words, optional leading icon, `Variant` (`VariantDefault`, `VariantSecondary`, `VariantOutline`, `VariantDestructive`, `VariantGhost`, `VariantLink`).

```templ
@badge.Badge(badge.Props{Variant: badge.VariantDestructive}) { @icon.CircleAlert() Degraded }
```

- Badge has no success or warning variant. For state, use `patterns.StatusBadge` (positive, warning, destructive, neutral: icon, word and `positive` / `warning` tokens on a 10% tint) and map your states to it once. Use `VariantSecondary` for plain categories (type, team) and `VariantDefault` for emphasis like a Production env.
- Status badges carry an icon or word, never colour alone.
- Version tags and IDs: `VariantOutline` with `font-mono`.
