#!/bin/bash
set -eu

# hash the content of the dist directory
# and write the sum in .dist.sha256
cd "$(dirname "$0")/.."
echo -ne $(find dist -type f | xargs cat | sha256sum) > .dist.sha256
