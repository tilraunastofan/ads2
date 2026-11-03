# ThemeToggle

A ghost icon button that switches light and dark. `patterns.ThemeToggle()` from `github.com/tilraunastofan/ads2/patterns` (imported from the ads2 module, see `docs/setup.md`), wired to `assets/js/theme.js`; already in `patterns.Shell`'s top bar.

**Consumer provides:** nothing but `theme.js` loaded as a plain (not deferred) `<script>` in `<head>`, which `patterns.Base` does.

```templ
<header class="flex h-14 items-center justify-between gap-2 border-b px-4">
    @sidebar.Trigger()
    @patterns.ThemeToggle("en") // page language: "Dark theme" / "Dökkt þema"
</header>
```

- First visit follows the OS (`prefers-color-scheme`); after a click the choice is stored in `localStorage` and wins.
- Dark mode is the `dark` class on `<html>`: it swaps the token variables and enables `dark:`. `data-theme` is not used.
- `theme.js` runs before first paint, so there is no flash of the wrong theme.
- Until the user clicks, the page follows OS theme changes live.
- It is a toggle button: the name stays "Dark theme" and `theme.js` keeps `aria-pressed` in step (also after boosted navigation).
- The theme sets `color-scheme`, so scrollbars and native date inputs follow it.
