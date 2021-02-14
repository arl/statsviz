#!/bin/sh
set -eu

cd ${GITHUB_WORKSPACE:-$(git rev-parse --show-toplevel)}

# Re-generate assets file
go generate

# Check that the generated assets only modifications are date-related.
if git diff --unified=1 assets_vfsdata.go | grep -E '^(\+|\-)\s.*$' | grep -v modTime; then
  # some diff content is not date-related, failing.
  echo "assets_vfsdata.go is not up to date"
  exit 1
fi

echo "Ok, assets are up to date"
