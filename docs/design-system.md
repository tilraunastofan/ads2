# Arion Design System 2

Arion Design System 2 dresses **shadcn-templ v2** (nova style) in Arion banki's palette for Go web apps built with **Go + templ + htmx 4 + Tailwind v4**. It adds no component library of its own. Every component is a shadcn-templ component (`github.com/axadrn/shadcn-templ/v2/components/...`), and the system only sets the theme variables they read and the rules for using them. Read **Setup** to wire it into a project and **htmx patterns** for the server-driven interactions.

## Principles

- **shadcn first.** Reach for a shadcn-templ component before writing markup. Customise with `Class` (merged by `utils.CN`) and `Attributes`, never by forking the component. If a component is missing, compose existing ones and add a pattern card here.
- **The server owns state.** Pages render complete HTML; htmx swaps fragments. Use a component's `Default*` props for initial state and re-render from the server instead of syncing client state.
- **Quiet, precise, legible.** Arion blue is the only loud colour. Everything else is navy ink on a cool, pale canvas, with hairline rings instead of shadows. Numbers line up.

## Content fundamentals

- **Voice:** plain, factual, calm. State what happened and what to do next: "Sync failed. Helm render failed for payments-api: values.yaml line 42." No exclamation marks, no emoji, no jokes in errors.
- **Person:** address the user as "you" only when needed. Buttons are imperative verbs ("Promote to prod", "Open pull request"). Never "Click here" or "OK".
- **Casing:** sentence case everywhere: titles, buttons, menu items, table headers. Product and system names keep their own casing (OpenShift, Argo CD, Kargo).
- **Language:** English for internal platform tools. Customer-facing and finance views follow Arion and use Icelandic ("Fjárhagsupplýsingar", "Hagnaður 3F 2026"). Never mix the two inside one sentence, which includes units: "41 min" on an English page, "41 mín." on an Icelandic one.
- **Numbers:** one format in every language, Arion's: `.` groups thousands and `,` is the decimal mark (`1.284.500`, `2,3%`, `+9,5%`), also on English pages. Only the words and units around the number follow the page language (`1.284.500 kr.` / `128,4 m.kr.`). Put units in the table caption or column header, not in every cell. Losses use a real minus sign `−`. Dates read `23. október 2026` in Icelandic UI and `23 Oct 2026` in English UI. Times use 24 hours. Format with `patterns.Int`, `Dec`, `Signed` and `Pct` rather than by hand.
- **Identifiers:** image tags, SHAs, namespaces and commands are set in `code` (JetBrains Mono).

## Color

Tokens are shadcn's own variable names, so shadcn-templ components pick them up unchanged. The light theme comes from arionbanki.is; the dark theme is derived from Arion's navy family.

- **Canvas and surfaces.** Pages sit on `background` (Arion's `#f7f9fa`). Cards, tables, dialogs and menus use `card` and `popover` (white), separated by the components' 10% `foreground` ring. Use `muted` for tracks, tab lists, skeletons and the pressed state.
- **Text.** Use `foreground` (Arion navy `#23313d`) for all body text. `muted-foreground` is for descriptions, helper text and captions only. Strong secondary labels use `slate`.
- **Action.** `primary` (Arion blue `#005ab4`) is reserved for the one primary button per view, links, checked controls, the active tab underline and progress. Do not use it for decoration or headings.
- **Hover and selection.** `accent` (pale Arion blue) with `accent-foreground` covers ghost and menu hover, the active sidebar item and highlighted rows. Text selection uses `selection` (`#acd5ff`, Arion's focus fill).
- **Status.** Errors use `destructive` (Arion red, darkened to `#cc0000`). Gains and healthy states use `positive` (teal, an addition to the Arion palette). In-progress, degraded and caution states use `warning` (amber-brown, also an addition). Teal and red stay apart for all common forms of colour blindness. Status is never shown by colour alone: always add a word, icon, sign or arrow. Use `patterns.StatusBadge` for state pills (Healthy / Progressing / Degraded) so every app maps states the same way; each token has a `*-foreground` for solid fills.
- **Editorial accents.** `ink` is Arion's dark CTA (one per page at most, never beside a blue primary). `gold` and `plum` are for highlight figures and charts. `sand` is a warm band behind report covers.
- **Charts.** Use `chart-1` to `chart-5` in order: blue, gold, plum, slate, light slate. Label series directly; the five are tuned to stay apart under red–green colour blindness.
- **Contrast.** Every text token clears 4.5:1, and every control border, focus ring and chart mark clears 3:1, on the surfaces its usage note names, in both themes. `cmd/gentheme` checks this on every `task theme` and `task check`, so a token edit that breaks a claim fails the build. Four Arion source values were adjusted to pass: `muted-foreground` `#5e7788` → `#536a7a`, `destructive` `#d70000` → `#cc0000`, the `input` border (50% navy, 2.9:1) → solid `#7b8c98`, and light `chart-5` → `#8294a0`. Dark `destructive` is `#ff7a70` so nova's 20% tint badge stays legible.

## Typography

- **Families.** The `sans` family (Hanken Grotesk) carries all UI and reading text. The `display` family (Source Serif 4) is for page titles and heroes only, mirroring Arion's serif headlines. The `mono` family (JetBrains Mono) is for figures in tables and for code.
- **Licensing.** The system ships open fonts only, self-hosted from `assets/fonts/` (SIL OFL), so no CDN is needed. Arion's own faces, Suisse Intl and Suisse Works (Swiss Typefaces, licensed), are an optional opt-in for teams that receive the files from Arion's UI team: add their `@font-face` rules to `fonts.css` and, after the `arion-theme.css` import in `globals.css`, put them first in the stacks (`:root { --font-sans: "Suisse Intl", "Hanken Grotesk", ui-sans-serif, system-ui, sans-serif; }`, and `--font-display: "Suisse Works", "Source Serif 4", ui-serif, Georgia, serif`). Nothing else changes.
- **Display styles.** One `display-*` title per page, at weight 600: `display-sm` (40px) on app pages, `display-lg` or `display-xl` on landing and report covers.
- **Utilities.** Every style in `tokens.json` is a Tailwind utility with its line height, weight and tracking: `text-display-sm`, `text-heading-1…4`, `text-body-lg`, `text-body`, `text-body-sm`, `text-label`, `text-caption`, `text-figure`, `text-code`. Use `font-display text-display-sm` for a page title and `font-mono text-code` for identifiers; never `text-[40px]`.
- **Headings.** `heading-1` to `heading-4` (sans, 600) are for sections inside a page. shadcn titles (`CardTitle`, `DialogTitle`) stay sans: `font-heading` is mapped to `--font-sans`.
- **Text.** Component text is shadcn's `text-sm` (`body-sm`, 14px). Docs and long descriptions use `body` (16px). `caption` (12px) is for timestamps and footnotes only.
- **Numbers.** Amounts use `figure`: `font-mono tabular-nums`, right-aligned.

## Spacing, radii, elevation

- **Spacing.** Use Tailwind's 4px scale; the `space-*` tokens name the steps used. Inside a card, gaps are 16px (`--card-spacing`). Between cards and blocks, use `gap-6` (24px). Above a page title, 48px.
- **Radii.** Follow Arion: 4px `radius-sm` for badges in tables, kbd and checkboxes; 8px `radius-lg` for buttons, inputs and alerts; 10px `radius-xl` for cards, dialogs and menus; `radius-full` for pills, avatars and switches. `--radius` is 0.5rem.
- **Borders over shadows.** Cards use a ring, not a shadow. The elevation tokens replace Tailwind's shadow scale, so the components pick them up without changes: `shadow-sm` is `elevation-1`, `shadow-md` (menus, selects) is `elevation-2`, `shadow-lg` (dialogs, toasts) is `elevation-3`. Overlay scrims stay shadcn's `bg-black/10`. Hairlines (`border`) divide rows and sections but never form a control's only boundary.
- **Focus.** Keep shadcn's focus treatment: a solid `border-ring` (6.7:1 on white) plus a 3px `ring/50` halo. Never remove outlines.

## Layout and motion

- **Internal tools.** Use `patterns.Shell`, the `sidebar` shell (`sidebar.Provider` → `Sidebar` + `Inset`) with the theme toggle. The inset holds a single column, max 1200px, padded 24px.
- **Public and report pages.** Use a top bar and a centered 880–1040px column.
- **Page skeleton.** Breadcrumb, then the display title with a status line and up to two actions, then KPI cards, then the main table or content (see **Patterns**: `PageHeader`, `Stat`, `DataCard`).
- **Motion.** Motion is shadcn's own: 100–200ms fades and zooms on overlays. Never animate data. `prefers-reduced-motion: reduce` removes all of it (the theme sets this globally); only spinners keep turning.

## Iconography

- **Icon set.** Use **Lucide**, through shadcn-templ's `components/icon` package (`@icon.Download()`, `@icon.CircleAlert()`). Stroke 2, size 16px inside components (they size `svg` themselves). Icons inherit `currentColor`.
- **Pairing with labels.** Icons accompany labels; icon-only buttons need `aria-label` and a `tooltip`. No emoji, no filled or duotone sets, no custom drawings.
- **Logos.** `assets/img/arion-mark.svg` is the star mark, single ink `#005ab4` (`primary`). It goes in the sidebar header and favicons, at 24–32px, on light surfaces or `card`. It is a vector trace of Arion's logo and is the version this system uses; never redraw or recolour it. The full lockup (charcoal wordmark) is for `background` or `card` in the light theme only; in dark mode use the mark alone. Give every logo clear space of at least half its height, and never stretch it or place it on `primary`.

## Components

These are shadcn-templ components, themed. Each card shows the markup the component renders, and each README gives the templ call and htmx usage.

- **Actions:** Button, ButtonGroup, Toggle, ToggleGroup
- **Forms:** Field, Label, Input, InputGroup, Textarea, Checkbox, RadioGroup, Switch, NativeSelect, Select, Slider
- **Navigation:** Sidebar, Breadcrumb, Tabs, Pagination, Accordion
- **Data display:** Card, Table, Badge, Avatar, Item, Kbd, Separator
- **Feedback:** Alert, Toast, Progress, Spinner, Skeleton, Empty
- **Overlays:** Dialog, AlertDialog, DropdownMenu, Tooltip
- **Patterns** (`github.com/tilraunastofan/ads2/patterns`, see `docs/patterns.md`): `Shell`, `PageHeader`, `Stat`, `StatusBadge`, `DataCard`, `ThemeToggle`, number formatting, and the Dashboard page composition

Every component above is installed in this repo, as are Sheet and AspectRatio. Other shadcn-templ components (Calendar, DatePicker, Combobox, Command, Popover, HoverCard, Drawer, ContextMenu, Menubar, NavigationMenu, Chart, Carousel, Resizable, ScrollArea, InputOTP, Collapsible) inherit the same theme with no extra rules. Use them as shadcn documents them.

## Dark mode

Put `class="dark"` on `<html>` (shadcn-templ's convention); that is what swaps the token variables and enables Tailwind's `dark:` variant (`data-theme` is not used). `assets/js/theme.js` and `patterns.ThemeToggle` start from the OS setting on first visit, then remember the user's choice. Every token has a dark value. In dark mode, blue fills take dark text (`primary-foreground` is `#0c171b`), so always use the `*-foreground` token on a fill, never literal white.
