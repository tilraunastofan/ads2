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
- **Language:** English for internal platform tools. Customer-facing and finance views follow Arion and use Icelandic ("Fjárhagsupplýsingar", "Hagnaður 3F 2026"). Never mix the two inside one sentence.
- **Numbers:** Icelandic format: `.` groups thousands and `,` is the decimal mark (`1.284.500 kr.`, `+9,5%`). Use `m.kr.` for millions and put units in the table caption or column header, not in every cell. Losses use a real minus sign `−`. Dates read `23. október 2026` in Icelandic UI and `23 Oct 2026` in English UI. Times use 24 hours.
- **Identifiers:** image tags, SHAs, namespaces and commands are set in `code` (JetBrains Mono).

## Color

Tokens are shadcn's own variable names, so shadcn-templ components pick them up unchanged. The light theme comes from arionbanki.is; the dark theme is derived from Arion's navy family.

- **Canvas and surfaces.** Pages sit on `background` (Arion's `#f7f9fa`). Cards, tables, dialogs and menus use `card` and `popover` (white), separated by the components' 10% `foreground` ring. Use `muted` for tracks, tab lists, skeletons and the pressed state.
- **Text.** Use `foreground` (Arion navy `#23313d`) for all body text. `muted-foreground` is for descriptions, helper text and captions only. Strong secondary labels use `slate`.
- **Action.** `primary` (Arion blue `#005ab4`) is reserved for the one primary button per view, links, checked controls, the active tab underline and progress. Do not use it for decoration or headings.
- **Hover and selection.** `accent` (pale Arion blue) with `accent-foreground` covers ghost and menu hover, the active sidebar item and highlighted rows. Text selection uses `selection` (`#acd5ff`, Arion's focus fill).
- **Status.** Errors use `destructive` (Arion red `#d70000`). Gains and healthy states use `positive` (teal, an addition to the Arion palette). Teal stays off the red–green axis so colour-blind readers can tell the two apart. Status is never shown by colour alone: always add a word, icon, sign or arrow.
- **Editorial accents.** `ink` is Arion's dark CTA (one per page at most, never beside a blue primary). `gold` and `plum` are for highlight figures and charts. `sand` is a warm band behind report covers.
- **Charts.** Use `chart-1` to `chart-5` in order: blue, gold, plum, slate, light slate. Label series directly.
- **Contrast.** Every text token clears 4.5:1 on the surfaces its usage note names, in both themes. Two source values sit just below their target and are kept exact: `input` at 2.9:1 on white, and `destructive` on its own 10% tint at 4.48:1. `muted-foreground` was darkened from Arion's `#5e7788` to `#536a7a` to pass.

## Typography

- **Families.** The `sans` family (Hanken Grotesk) carries all UI and reading text. The `display` family (Source Serif 4) is for page titles and heroes only, mirroring Arion's serif headlines. The `mono` family (JetBrains Mono) is for figures in tables and for code.
- **Licensing.** Arion's own faces are Suisse Intl and Suisse Works (Swiss Typefaces, licensed). Both are listed second in the stacks: if you hold a licence and self-host them, they take over automatically. Otherwise the open substitutes render.
- **Display styles.** One `display-*` title per page, at weight 600: `display-sm` (40px) on app pages, `display-lg` or `display-xl` on landing and report covers.
- **Headings.** `heading-1` to `heading-4` (sans, 600) are for sections inside a page. shadcn titles (`CardTitle`, `DialogTitle`) stay sans: `font-heading` is mapped to `--font-sans`.
- **Text.** Component text is shadcn's `text-sm` (`body-sm`, 14px). Docs and long descriptions use `body` (16px). `caption` (12px) is for timestamps and footnotes only.
- **Numbers.** Amounts use `figure`: `font-mono tabular-nums`, right-aligned.

## Spacing, radii, elevation

- **Spacing.** Use Tailwind's 4px scale; the `space-*` tokens name the steps used. Inside a card, gaps are 16px (`--card-spacing`). Between cards and blocks, use `gap-6` (24px). Above a page title, 48px.
- **Radii.** Follow Arion: 4px `radius-sm` for badges in tables, kbd and checkboxes; 8px `radius-lg` for buttons, inputs and alerts; 10px `radius-xl` for cards, dialogs and menus; `radius-full` for pills, avatars and switches. `--radius` is 0.5rem.
- **Borders over shadows.** Cards use a ring, not a shadow. `elevation-2` is for floating menus and `elevation-3` for dialogs and toasts only. Hairlines (`border`) divide rows and sections but never form a control's only boundary.
- **Focus.** Keep shadcn's focus treatment: a solid `border-ring` (6.7:1 on white) plus a 3px `ring/50` halo. Never remove outlines.

## Layout and motion

- **Internal tools.** Use the `sidebar` shell (`sidebar.Provider` → `Sidebar` + `Inset`). The inset holds a single column, max 1200px, padded 24px.
- **Public and report pages.** Use a top bar and a centered 880–1040px column.
- **Page skeleton.** Breadcrumb, then the display title with a status line and up to two actions, then KPI cards, then the main table or content (see the **Dashboard** pattern).
- **Motion.** Motion is shadcn's own: 100–200ms fades and zooms on overlays. Never animate data. Respect `prefers-reduced-motion`.

## Iconography

- **Icon set.** Use **Lucide**, through shadcn-templ's `components/icon` package (`@icon.Download()`, `@icon.CircleAlert()`). Stroke 2, size 16px inside components (they size `svg` themselves). Icons inherit `currentColor`.
- **Pairing with labels.** Icons accompany labels; icon-only buttons need `aria-label` and a `tooltip`. No emoji, no filled or duotone sets, no custom drawings.
- **Logos.** Use the uploaded files and never redraw them. `Arion SVG/logo.svg` is the star mark, single ink `#005ab4` (`primary`). It goes in the sidebar header and favicons, at 24–32px, on light surfaces or `card`. `Arion logos/arion-banki-logo-txt.webp` is the full lockup with a charcoal wordmark: use it on `background` or `card` in the light theme only, and in dark mode use the mark alone. `Arion logos/arion-logo.png` is a raster mark for places SVG can't go (email, PDF). Give every logo clear space of at least half its height, and never recolour, stretch or place it on `primary`.

## Components

These are shadcn-templ components, themed. Each card shows the markup the component renders, and each README gives the templ call and htmx usage.

- **Actions:** Button, ButtonGroup, Toggle, ToggleGroup
- **Forms:** Field, Label, Input, InputGroup, Textarea, Checkbox, RadioGroup, Switch, NativeSelect, Select, Slider
- **Navigation:** Sidebar, Breadcrumb, Tabs, Pagination, Accordion
- **Data display:** Card, Table, Badge, Avatar, Item, Kbd, Separator
- **Feedback:** Alert, Toast, Progress, Spinner, Skeleton, Empty
- **Overlays:** Dialog, AlertDialog, DropdownMenu, Tooltip
- **Patterns:** Dashboard

Other shadcn-templ components (Calendar, DatePicker, Combobox, Command, Popover, HoverCard, Sheet, Drawer, ContextMenu, Menubar, NavigationMenu, Chart, Carousel, Resizable, ScrollArea, InputOTP, Collapsible, AspectRatio) inherit the same theme with no extra rules. Use them as shadcn documents them.

## Dark mode

Put `class="dark"` on `<html>` (shadcn-templ's convention). The `[data-theme="dark"]` attribute is supported too. Every token has a dark value. In dark mode, blue fills take dark text (`primary-foreground` is `#0c171b`), so always use the `*-foreground` token on a fill, never literal white.
