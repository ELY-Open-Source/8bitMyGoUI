#!/bin/sh
# Uploads a version, tests its preview URL, then promotes it.
set -eu
cd "$(dirname "$0")"
if [ -n "$(git status --porcelain)" ]; then
  echo "deploy: commit first; each version names its commit" >&2
  exit 1
fi
./build.sh
rev=$(git rev-parse --short HEAD)
out=$(npx wrangler versions upload --message "$rev")
echo "$out"
id=$(echo "$out" | sed -n 's/^Worker Version ID: //p')
preview=$(echo "$out" | sed -n 's/^Version Preview URL: //p')
if [ -z "$id" ] || [ -z "$preview" ]; then
  echo "deploy: wrangler printed no version id or preview URL" >&2
  exit 1
fi
BASE_URL=$preview npx playwright test
npx wrangler versions deploy "$id@100%" --yes --message "$rev"
site=$(echo "$preview" | sed 's#https://[0-9a-f]*-#https://#')
main=$(sed -n 's#.*src="\(a/main\.[0-9a-f]*\.js\)".*#\1#p' dist/index.html)
tries=0
until curl -fs "$site/" | grep -q "$main"; do
  tries=$((tries + 1))
  if [ "$tries" -ge 30 ]; then
    echo "deploy: $site still serves an older build than $main" >&2
    exit 1
  fi
  sleep 2
done
# A first fetch warms the edge before the browser tests.
wasm=$(sed -n "s#.*fetch('\\(a/app\\.[0-9a-f]*\\.wasm\\)').*#\\1#p" dist/a/main.*.js)
curl -fs --max-time 120 -o /dev/null "$site/$wasm"
BASE_URL=$site npx playwright test
echo "deploy: $rev live at $site (version $id)"
