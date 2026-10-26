# ThemeToggle

A ghost icon button that switches light and dark. `patterns.ThemeToggle()` from `github.com/tilraunastofan/ads2/patterns` (copy the package into your app like `components/`), wired to `assets/js/theme.js`; already in `patterns.Shell`'s top bar.

**Consumer provides:** nothing but `theme.js` loaded as a plain (not deferred) `<script>` in `<head>`, which `patterns.Base` does.

```templ
<header class="flex h-14 items-center justify-between gap-2 border-b px-4">
    @sidebar.Trigger()
    @patterns.ThemeToggle()
</header>
```

- First visit follows the OS (`prefers-color-scheme`); after a click the choice is stored in `localStorage` and wins.
- Dark mode is the `dark` class on `<html>`: it swaps the token variables and enables `dark:`. `data-theme` is not used.
- `theme.js` runs before first paint, so there is no flash of the wrong theme.
