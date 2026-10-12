#!/bin/sh
# Ships files first, then the page, so none goes missing.
set -eu
cd "$(dirname "$0")"
site=https://8bitmygoui.zhangyanghaha0407.workers.dev
if [ -n "$(git status --porcelain)" ]; then
  echo "deploy: commit first; each version names its commit" >&2
  exit 1
fi
# patient runs a check every 2 s for a minute.
patient() {
  tries=0
  until "$@"; do
    tries=$((tries + 1))
    if [ "$tries" -ge 30 ]; then
      echo "deploy: gave up waiting for: $*" >&2
      exit 1
    fi
    sleep 2
  done
}
# upload makes a version from the assets in $1.
upload() {
  out=$(npx wrangler versions upload --assets "$1" --message "$2")
  echo "$out" >&2
  id=$(echo "$out" | sed -n 's/^Worker Version ID: //p')
  preview=$(echo "$out" | sed -n 's/^Version Preview URL: //p')
  if [ -z "$id" ] || [ -z "$preview" ]; then
    echo "deploy: wrangler printed no version id or preview URL" >&2
    exit 1
  fi
}
./build.sh
rev=$(git rev-parse --short HEAD)
main=$(sed -n 's#.*src="\(a/main\.[0-9a-f]*\.js\)".*#\1#p' dist/index.html)
stage=$(mktemp -d)
trap 'rm -rf "$stage"' EXIT
# The live build's files ride along, for pages already open.
if ! curl -fs "$site/" -o "$stage/index.html"; then
  echo "deploy: cannot read the live page at $site" >&2
  exit 1
fi
for f in $(sed -n 's#.*src="\(a/[^"]*\)".*#\1#p' "$stage/index.html"); do
  curl -fs "$site/$f" -o "dist/$f"
done
old=$(sed -n 's#.*src="\(a/main\.[0-9a-f]*\.js\)".*#\1#p' "$stage/index.html")
wasm=$(sed -n "s#.*fetch('\\(a/app\\.[0-9a-f]*\\.wasm\\)').*#\\1#p" "dist/$old")
curl -fs --max-time 300 "$site/$wasm" -o "dist/$wasm"
cp -R dist/a dist/_headers "$stage/"
# Step one: the live page beside the new files.
upload "$stage" "$rev files"
npx wrangler versions deploy "$id@100%" --yes --message "$rev files"
patient curl -fs -o /dev/null "$site/$main"
# Step two: the new page, tested on its preview first.
upload ./dist "$rev"
BASE_URL=$preview npx playwright test
npx wrangler versions deploy "$id@100%" --yes --message "$rev"
patient sh -c "curl -fs '$site/' | grep -q '$main'"
BASE_URL=$site npx playwright test
echo "deploy: $rev live at $site (version $id)"
