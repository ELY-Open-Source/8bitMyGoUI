#!/bin/sh
# Builds web/dist: the gallery as wasm plus the page shell.
set -eu
cd "$(dirname "$0")/.."
build=$(git rev-parse --short HEAD)-$(date +%s)
rm -rf web/dist
mkdir -p web/dist
GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o web/dist/app.wasm ./cmd/web
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/dist/
for f in index.html main.js; do
  sed "s/__BUILD__/$build/g" "web/$f" > "web/dist/$f"
done
echo "web/dist built: $build, app.wasm $(wc -c < web/dist/app.wasm) bytes"
