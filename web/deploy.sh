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
BASE_URL=$site npx playwright test
echo "deploy: $rev live at $site (version $id)"
