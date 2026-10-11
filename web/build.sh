#!/bin/sh
# Builds web/dist; files under a/ are named by content.
set -eu
cd "$(dirname "$0")/.."
d=web/dist
rm -rf $d
mkdir -p $d/a
GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o $d/app.wasm ./cmd/web
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" $d/
size=$(wc -c < $d/app.wasm | tr -d ' ')
sed "s/__SIZE__/$size/g" web/main.js > $d/main.js
# name moves a file into a/ by its content hash.
name() {
  to=a/$1.$(git hash-object "$d/$2" | cut -c1-12).$3
  mv "$d/$2" "$d/$to"
  echo "$to"
}
wasm=$(name app app.wasm wasm)
sed -i.bak "s#__WASM__#$wasm#g" $d/main.js && rm $d/main.js.bak
exec=$(name wasm_exec wasm_exec.js js)
main=$(name main main.js js)
sed "s#__EXEC__#$exec#g; s#__MAIN__#$main#g" web/index.html > $d/index.html
cp web/_headers $d/
echo "web/dist built: $main $wasm ($size bytes)"
