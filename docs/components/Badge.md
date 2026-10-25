# Badge

Short status or metadata label: `badge.Badge` from `components/badge` (`<span>`, or `<a>` with `Href`).

**Consumer provides:** one to three words, optional leading icon, `Variant` (`VariantDefault`, `VariantSecondary`, `VariantOutline`, `VariantDestructive`, `VariantGhost`, `VariantLink`).

```templ
@badge.Badge(badge.Props{Variant: badge.VariantDestructive}) { @icon.CircleAlert() Degraded }
```

- Map health to variants once, in a helper (`StatusBadge(h Health)`): Healthy → secondary, Progressing → outline, Degraded → destructive, Production env → default.
- Status badges carry an icon or word, never colour alone.
- Version tags and IDs: `VariantOutline` with `font-mono`.
