# Migrating an app from copied files to the ads2 module

Until `v0.1.0`, apps started from ads2 by **copying** `components/`, `patterns/`, `utils/` and the assets (theme CSS, fonts, htmx, helper scripts) into their own repository. Fixes in ads2 never reached them. From `v0.1.0` apps **import** ads2 as a Go module instead (see **How apps use ads2** in `docs/setup.md`), and `go get -u` brings every update.

This guide converts an app that copied the files. It is written to be followed by an agent, step by step; each step says how to check it before going on. Do the whole migration on a branch, as one pull request.

## 0. Is this app a copy?

The app copied ads2 if **any** of these hold. Check them all and note which:

- `grep -rl "Arion Design System\|arion-theme.css" --include=*.go --include=*.templ --include=*.css .` finds files.
- There is a `patterns/` package with `Shell`, `PageHeader`, `Stat` (`grep -l "templ Shell(" patterns/*.templ`).
- There is a `components/` folder with shadcn-templ components and a `components.json` with `"style": "base-nova"`.
- `assets/css/arion-theme.css`, `assets/css/fonts.css`, `assets/js/toast-guard.js` or `assets/js/theme.js` exist.

If none hold, the app doesn't use ads2: stop. If it already imports `github.com/tilraunastofan/ads2/patterns`, it is migrated: go to **Updating later**.

Read the app's `go.mod` for its module path (`APP` below, for example `example.com/payments-ui`) and Go version (ads2 needs Go 1.26+; raise it if lower).

## 1. Find local changes before deleting anything

The copies may have been edited. Deleting an edited file loses the change, so compare every copied file with ads2 first.

```sh
go mod download -json github.com/tilraunastofan/ads2@v0.1.0   # prints "Dir": the module's files
ADS2=<that Dir>
for d in components patterns utils; do
  diff -r -x '*_templ.go' "$d" "$ADS2/$d"
done
for f in assets/css/arion-theme.css assets/css/fonts.css assets/js/theme.js assets/js/toast-guard.js \
         assets/js/filter-form.js assets/js/request-errors.js assets/js/dialog-form.js \
         assets/js/datepicker.js assets/js/focus-after-swap.js assets/js/htmx.min.js; do
  [ -f "$f" ] && ! cmp -s "$f" "$ADS2/$f" && echo "differs: $f"
done
```

Sort every difference into one of three kinds, and write the list into the pull request description:

| Difference | What to do |
| --- | --- |
| ads2 is **newer** (the app copied an older ads2; `git log` in ads2 shows the change) | Nothing: the module brings it. |
| An **import path** only (`APP/components/…` instead of `github.com/tilraunastofan/ads2/components/…`) | Nothing: step 3 rewrites them. |
| A **local change** the app made (a new prop, a different class, an extra pattern) | Don't drop it. If it is app-specific, move it into the app's own package (for example `internal/ui`), built on the module's components. If every app would want it, stop and propose it to ads2 first (an issue or a PR there), then migrate onto the ads2 version that has it. Never keep a modified copy of an ads2 package: it would load scripts twice and silently miss updates. |

A component in `components/` that ads2 **doesn't** have (`diff` says "Only in components: …") is app-specific: either propose it to ads2, or keep it and see **Components ads2 doesn't have** below.

Check: every difference is in the list with a decision.

## 2. Add the module

```sh
go get github.com/tilraunastofan/ads2@v0.1.0
go get -tool github.com/tilraunastofan/ads2/cmd/ads2css@v0.1.0
go get -tool github.com/a-h/templ/cmd/templ@v0.3.1100   # ads2's templ version (see its go.mod)
```

Check: `go list -m github.com/tilraunastofan/ads2` prints `v0.1.0`.

## 3. Point the imports at the module

Rewrite the import paths in every `.go` and `.templ` file the app owns (not in the folders you are about to delete):

```sh
APP=$(go list -m)
grep -rl --include=*.go --include=*.templ -e "\"$APP/components" -e "\"$APP/patterns" -e "\"$APP/utils" . \
  | grep -v -e '^./components/' -e '^./patterns/' -e '^./utils/' \
  | xargs sed -i.bak \
      -e "s#\"$APP/components/#\"github.com/tilraunastofan/ads2/components/#g" \
      -e "s#\"$APP/components\"#\"github.com/tilraunastofan/ads2/components\"#g" \
      -e "s#\"$APP/patterns\"#\"github.com/tilraunastofan/ads2/patterns\"#g" \
      -e "s#\"$APP/utils\"#\"github.com/tilraunastofan/ads2/utils\"#g"
find . -name '*.bak' -delete
```

If step 1 moved local changes into an app package, point those call sites at it.

Check: `grep -rn "\"$APP/\(components\|patterns\|utils\)" --include=*.go --include=*.templ . | grep -v '^./\(components\|patterns\|utils\)/'` prints nothing.

## 4. Delete the copies

Delete only what step 1 cleared:

```sh
git rm -r components patterns utils components.json
git rm assets/css/arion-theme.css assets/css/fonts.css assets/css/tw-animate.css assets/css/shadcn-tailwind.css
git rm -r assets/fonts
git rm assets/js/theme.js assets/js/toast-guard.js assets/js/filter-form.js assets/js/request-errors.js \
       assets/js/dialog-form.js assets/js/datepicker.js assets/js/focus-after-swap.js assets/js/htmx.min.js \
       assets/js/shadcn-templ-*.js assets/img/arion-mark.svg assets/img/arion-mark-favicon.svg
```

`git rm` fails on a file that isn't there; drop it from the command (older copies lack the newer scripts). Keep everything else in `assets/` that is the app's own: its images, its own scripts, `globals.css`.

Check: `go build ./...` fails only on the asset serving (next step) and on API differences (step 7), never on a missing package.

## 5. Serve the assets from the module

Find where the app serves `/assets/` (`grep -rn '"/assets/\|/assets/"' --include=*.go .`), usually an `embed.FS` with a `http.FileServer`. Replace it with the module's handler, the app's own files on top:

```go
import "github.com/tilraunastofan/ads2/assets"

//go:embed static
var static embed.FS // the app's own files: css/output.css, its images

app, _ := fs.Sub(static, "static")
mux.Handle("GET /assets/", http.StripPrefix("/assets/", assets.Handler(assets.Overlay(app))))
```

Keep the app's embed root as it is if it already holds only app files; `fs.Sub` it so that `css/output.css` sits at the root of `app`. If the app copied the showcase's own asset handler (ETags, gzip, a `SHADCN_TEMPL_DEV` disk mode), delete it: `assets.Handler` does the same. For a dev mode that reads the app's files from disk, overlay `os.DirFS("static")` instead of the embed.

Wrap the handler in `patterns.SidebarState(…)` if the app uses `patterns.Shell` and doesn't yet.

Check: `go build ./...` passes, or fails only in step 7's API differences.

## 6. Connect Tailwind to the module

Replace the app's `globals.css` with:

```css
@import "tailwindcss";
@import "./ads2.gen.css";
```

Keep any app-specific CSS after these lines (custom `@layer components`, `@source` lines for app folders Tailwind doesn't find). Remove what ads2 now provides: imports of `arion-theme.css`, `tw-animate.css`, `shadcn-tailwind.css`, the `@layer base` with `border-border`.

Add `ads2css` before every Tailwind run (Taskfile, Makefile, CI, Dockerfile), and gitignore its output and Tailwind's:

```sh
go tool ads2css -o <folder of globals.css>/ads2.gen.css
tailwindcss -i <globals.css> -o <output.css> --minify
printf '%s\n' '<folder>/ads2.gen.css' '<folder>/output.css' >> .gitignore
```

Remove `go tool shadcn-templ bundle` and the `shadcn-templ` tool from the build and `go.mod` (`go get -tool github.com/axadrn/shadcn-templ/v2/cmd/shadcn-templ@none`): the bundle comes from the module now.

Check: after the build, `output.css` contains `.text-display-sm` and `.bg-primary` (`grep -c 'text-display-sm' <output.css>`).

## 7. Fix API differences

```sh
go tool templ generate && go mod tidy && go build ./... && go vet ./...
```

The compiler lists every call that changed between the app's copy and the module. Read `CHANGELOG.md` in ads2 (in `$ADS2`) for each entry since the app copied, which says what changed and how to update a call. Fix each, regenerate, rebuild, until it passes. Don't work around a change by copying the old code back.

Check: `go build ./... && go vet ./... && go test ./...` pass.

## 8. Verify in a browser

Run the app and open each page in Chrome (a dark-mode extension such as Dark Reader repaints pages; disable it when judging colours):

- No `404` for anything under `/assets/` in the network panel, no CSP errors in the console.
- Fonts: the page title is serif (Source Serif 4), body text Hanken Grotesk.
- Dark mode: the theme toggle switches it, and the Arion mark turns white.
- The sidebar collapses to icons and stays collapsed on the next page.
- A dialog or sheet opens and closes; a toast appears and the tab stays responsive while you follow a link during it.
- Numbers show in Arion format (`1.284.500`), also in chart tooltips.

## 9. Commit

One pull request with the decisions from step 1 in its description. Note in the app's README or CLAUDE.md that the design system comes from `github.com/tilraunastofan/ads2` and is updated with `go get -u`.

## Updating later

```sh
go get -u github.com/tilraunastofan/ads2 && go mod tidy
go get -tool github.com/tilraunastofan/ads2/cmd/ads2css   # same version as the module
```

Read the `CHANGELOG.md` entries between the old and new version, apply the steps any of them name, then run the build (ads2css, Tailwind, templ), the tests, and step 8.

## Components ads2 doesn't have

The module ships every component in ads2's `components/`. An app that needs another one should get it added to ads2 (`shadcn-templ add <name>` there, plus a `docs/components/` page). If that can't wait, keep only that component in the app, under its own package path, and load its script yourself; but its script must not be one ads2's bundle already contains, or it runs twice. Move it to ads2 as soon as possible.
