# Fonts

Self-hosted variable fonts (latin subset), loaded by `assets/css/fonts.css`.

| File | Family | Source | Licence |
| --- | --- | --- | --- |
| `hanken-grotesk-latin-wght-normal.woff2` | Hanken Grotesk | `@fontsource-variable/hanken-grotesk` | SIL OFL 1.1 |
| `source-serif-4-latin-wght-normal.woff2` | Source Serif 4 | `@fontsource-variable/source-serif-4` | SIL OFL 1.1 |
| `jetbrains-mono-latin-wght-normal.woff2` | JetBrains Mono | `@fontsource-variable/jetbrains-mono` | SIL OFL 1.1 |

To refresh, download `files/<name>-latin-wght-normal.woff2` from the package on jsDelivr.
Arion's Suisse Intl and Suisse Works are licensed by Swiss Typefaces, not included, and optional: see docs/design-system.md (Typography).

## Caching and updates

The showcase serves this folder with `Cache-Control: public, max-age=31536000, immutable`, so browsers never re-download a font. **If you replace a font file, give it a new file name** (e.g. add a version suffix) and update `fonts.css` and the preload links in `patterns.Base`; otherwise browsers keep the old file for a year.

The `* Fallback` faces in `fonts.css` are metric-matched to these files. If you change a font, recompute `size-adjust`, `ascent-override` and `descent-override` (letter-weighted average advance width and hhea metrics of the new font over Arial, Times New Roman or Courier New).
