#!/usr/bin/env bash
# fetch-specs.sh — (re)download the pristine spec HTML into docs/specs/html.
# Normally never needed; the HTML is committed. Run `make specs` after.
set -euo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/docs/specs/html"
mkdir -p "$DIR"

fetch() { echo "fetching $2"; curl -sSL -o "$DIR/$1" "$2"; }

fetch xmlschema11-1.html        https://www.w3.org/TR/xmlschema11-1/
fetch xmlschema11-2.html        https://www.w3.org/TR/xmlschema11-2/
fetch xpath20.html              https://www.w3.org/TR/xpath20/
fetch xsd-precisionDecimal.html https://www.w3.org/TR/xsd-precisionDecimal/
