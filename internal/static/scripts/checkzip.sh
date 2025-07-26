#!/bin/bash
set -eu

# Check that the files in dist.zip are up to date with the files in dist/
# This must be run after npm run build.
cd "$(dirname "$0")/.."

if unzip -l dist.zip | grep -oP 'dist/assets/index-.*.(js|css)' | xargs ls >/dev/null; then
  echo "'dist.zip' contains up to date files."
else
  echo "'dist.zip' does not contain up to date files, did you forget to build/commit the assets?" && exit 1
fi
