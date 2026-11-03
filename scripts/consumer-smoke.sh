#!/usr/bin/env bash
# Builds testdata/consumer, the smallest app on the ads2 module, against this
# checkout (its tracked and new files, as the published module has them) and checks that it serves its page and every asset the page links to.
# This is what an app does (docs/setup.md, "New project"), so a change that breaks
# apps importing ads2 fails here. Needs go, tailwindcss and curl; uses the network
# the first time for the app's dependencies.
set -euo pipefail

root=$(cd "$(dirname "$0")/.." && pwd)
work=$(mktemp -d)
port=${PORT:-8099}
# wait reports the killed server's status (143); under set -e that would become
# the script's exit status, so the cleanup swallows it.
trap '{ kill "${pid:-}" && wait "${pid:-}"; } 2>/dev/null || true; rm -rf "$work"' EXIT

# The module as an app downloads it: tracked and new files, never ignored ones. A
# file the repo ignores (a generated bundle, say) is missing from the published
# module, so it must be missing here too.
mod="$work/ads2"
mkdir "$mod"
git -C "$root" ls-files -z --cached --others --exclude-standard |
  while IFS= read -r -d '' f; do [ -e "$root/$f" ] && printf '%s\0' "$f"; done |
  tar -C "$root" --null -T - -cf - | tar -C "$mod" -xf -

mkdir "$work/app"
cp -R "$root/testdata/consumer/." "$work/app/"
cd "$work/app"
templ_version=$(cd "$root" && go list -m -f '{{.Version}}' github.com/a-h/templ)
cat > go.mod <<GOMOD
module example.com/consumer

go 1.26

require github.com/tilraunastofan/ads2 v0.0.0

replace github.com/tilraunastofan/ads2 => $mod
GOMOD
go get -tool "github.com/a-h/templ/cmd/templ@$templ_version"
go get -tool github.com/tilraunastofan/ads2/cmd/ads2css
go tool templ generate
go mod tidy
go tool ads2css -o static/css/ads2.gen.css
tailwindcss -i static/css/globals.css -o static/css/output.css --minify
# Classes used only inside the module's patterns must be in the app's CSS.
for class in '.text-display-sm' 'collapsible\=icon'; do
  grep -qF -- "$class" static/css/output.css || { echo "output.css lacks .$class: ads2css @source lines missing?"; exit 1; }
done
go vet ./...
go build -o consumer .

PORT=$port ./consumer & pid=$!
for _ in $(seq 1 50); do curl -fs -o /dev/null "localhost:$port/" && break; sleep 0.1; done
page=$(curl -fs "localhost:$port/")
grep -q 'Built on the ads2 module.' <<<"$page"
# Every asset the page links to, from the module or the app's overlay.
links=$(grep -o '/assets/[^"]*' <<<"$page" | sort -u)
[ -n "$links" ]
for p in $links; do
  code=$(curl -s -o /dev/null -w '%{http_code}' "localhost:$port$p")
  [ "$code" = 200 ] || { echo "consumer: $p answered $code"; exit 1; }
done
echo "consumer smoke: ok ($(wc -l <<<"$links" | tr -d ' ') assets)"
