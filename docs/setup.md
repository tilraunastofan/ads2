# Setup

How to start a Go project on this system. The reference starter lives in the `ads2` repository; its `assets/css/arion-theme.css` is generated from this system's `tokens.json`.

## Stack

| Layer | Choice | Notes |
| --- | --- | --- |
| Server | Go 1.26+, `net/http` `ServeMux` | Handlers render templ components; no JSON for UI routes. |
| Templates | templ (`github.com/a-h/templ`) | Pin as a Go tool: `go get -tool github.com/a-h/templ/cmd/templ`. Commit generated `*_templ.go`; CI runs `templ generate` + `git diff --exit-code`. |
| Components | shadcn-templ v2 (`github.com/axadrn/shadcn-templ/v2`) | Style **nova** (`base-nova` in `components.json`). Add components with `shadcn-templ add <name>`. |
| Interactivity | htmx 4 | Vendored into `assets/js/` (on-prem: no runtime CDN). See **htmx patterns**. |
| CSS | Tailwind CSS v4.1+ standalone CLI | `assets/css/globals.css` → `assets/css/output.css`. |
| Fonts | Hanken Grotesk, Source Serif 4, JetBrains Mono | Self-host the woff2 files under `assets/fonts/` for air-gapped clusters; Google Fonts only for local dev. |

## New project

```sh
go install github.com/axadrn/shadcn-templ/v2/cmd/shadcn-templ@latest
shadcn-templ init my-app -t templ --preset nova
cd my-app
shadcn-templ add button card badge table field input input-group textarea label checkbox radio-group switch native-select select slider toggle toggle-group tabs breadcrumb pagination accordion sidebar avatar item kbd separator alert toast progress spinner skeleton empty dialog alert-dialog dropdown-menu tooltip icon
go mod tidy && go tool templ generate && go build ./... && go vet ./...
```

- Registry names are **kebab-case** (`alert-dialog`, `dropdown-menu`, `input-group`, `native-select`); the Go package names drop the hyphen (`alertdialog`, `dropdownmenu`…). The no-hyphen names return 404.
- `input-group` needs `textarea` installed too.
- In a fresh scaffold run `go mod tidy` first: `go tool templ` is declared in `go.mod` but only resolves after tidy. After that, run `go tool templ generate && go mod tidy` after every `.templ` edit or `shadcn-templ add`.
- Never edit the vendored `components/` or the generated `arion-theme.css`: change `tokens.json` and regenerate (`task theme`).

Then apply the Arion theme:

1. Copy `arion-theme.css` from `ads2/assets/css/` next to `globals.css`.
2. In `globals.css` delete the `:root { … }`, `.dark { … }`, `@theme inline { … }` and `@custom-variant dark` blocks that `init` wrote, and add `@import "./arion-theme.css";` after the other `@import` lines.
3. Keep `body` class `style-nova` (the init scaffold sets it). Keep menus solid: do not enable the translucent menu colour.
4. Load the fonts (self-hosted `@font-face` or the Google Fonts link in dev).

## Theme variables

`arion-theme.css` defines exactly shadcn's variables, so every shadcn-templ component is themed without changes, plus these extra `--color-*` utilities: `ink`, `ink-foreground`, `ink-hover`, `slate`, `gold`, `plum`, `sand`, `positive`, and `font-display` for the serif titles. Radii are Arion's 4 / 6 / 8 / 10px (`rounded-sm` … `rounded-xl`).

## Layout skeleton

```templ
package layouts

import (
    "my-app/components"
    "my-app/components/toast"
)

templ App(title string) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="utf-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1"/>
            <title>{ title } · Platform</title>
            <link rel="stylesheet" href="/assets/css/output.css"/>
            <script src="/assets/js/htmx.min.js" defer></script>
            @components.Scripts()
        </head>
        <body class="style-nova min-h-svh">
            { children... }
            @toast.Toaster()
        </body>
    </html>
}
```

## Verifying UI

- Dark mode is opt-in by class: `class="dark"` on `<html>` swaps the token variables. The theme ignores `prefers-color-scheme`; add your own toggle that sets the class. (`data-theme="dark"` only enables Tailwind's `dark:` variant; it does not change the tokens.)
- Check toasts, dialogs, menus and themes in a real browser; HTTP checks only cover the server side.
- Browser extensions such as Dark Reader repaint pages. Judge the theme by its variables (`getComputedStyle(document.documentElement).getPropertyValue('--background')`), not by screenshot colour.
- `task dev` needs ports 8090 and 7331 free (`lsof -nP -iTCP:8090 -sTCP:LISTEN`).

## On-prem checklist

- Vendor htmx (`assets/js/htmx.min.js`, 4.0.x) and the fonts. Serve `/assets/` from `go:embed` in the binary.
- Hashed `assets/js/shadcn-templ-*.js` bundles: `Cache-Control: immutable`. Everything else revalidates.
- Container: build with `task build`, run as non-root on OpenShift (random UID): write nothing to the image filesystem at runtime.
- CSP: the component scripts are files, so `script-src 'self'` works. If you use `hx-on:*`, add `htmx.config.inlineScriptNonce`.
