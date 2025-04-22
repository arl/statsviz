#!/bin/bash
set -eu

cd "$(dirname "$0")/.."
find dist -type f | xargs cat | sha256sum --check .dist.sha256 --quiet || \
  (echo "'dist' hash doesn not match the directory content, did you forget to build/commit the assets?" && exit 1)

