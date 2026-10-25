# ads2 — Arion Design System 2

A design system for Go web projects built with **Go + templ + htmx 4 + [shadcn-templ](https://shadcn-templ.com)**. It themes shadcn-templ v2 (nova style) with Arion banki's colours and adds usage rules, htmx patterns and a reference app.

- **Brand book and rules:** [`docs/design-system.md`](docs/design-system.md)
- **Wiring it into a project:** [`docs/setup.md`](docs/setup.md)
- **htmx 4 recipes:** [`docs/htmx-patterns.md`](docs/htmx-patterns.md)
- **Per-component guidance (templ + htmx):** [`docs/components/`](docs/components)
- **Tokens:** [`tokens/tokens.json`](tokens/tokens.json) → `go run ./cmd/gentheme` → [`assets/css/arion-theme.css`](assets/css/arion-theme.css)

## Use the theme in another project

```sh
shadcn-templ init --preset nova
cp /path/to/ads2/assets/css/arion-theme.css assets/css/
# in assets/css/globals.css: remove the :root / .dark / @theme inline blocks, then
#   @import "./arion-theme.css";
```

## Palette (light / dark)

| Token | Light | Dark | Use |
| --- | --- | --- | --- |
| `primary` | `#005ab4` | `#4da3ff` | The one action colour |
| `foreground` | `#23313d` | `#e6ecef` | Text |
| `background` | `#f7f9fa` | `#0c171b` | Canvas |
| `card` / `popover` | `#ffffff` | `#17212a` | Surfaces |
| `accent` | `#e0ebf6` | `#223c52` | Hover, active nav |
| `destructive` | `#d70000` | `#ff6b61` | Errors, losses |
| `positive` | `#00718a` | `#4fc1d6` | Gains, healthy |
| `ink` | `#17212a` | `#e6ecef` | Editorial dark CTA |

Fonts: Hanken Grotesk (UI), Source Serif 4 (page titles), JetBrains Mono (figures, code). Open substitutes for Arion's licensed Suisse Intl / Suisse Works.

## Showcase

```sh
task dev   # http://localhost:7331/apps
```

See `CLAUDE.md` for first-time bootstrap steps.
