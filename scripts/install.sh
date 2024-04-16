#!/usr/bin/env bash
 set -euo pipefail

echo "🏗️ Building Go binary"
go build mwi-redeploy.go;
echo "🚚 Moving binary to /use/local/bin";
cp mwi-redeploy /usr/local/bin/mwi-redeploy;
echo "✅ Done";
