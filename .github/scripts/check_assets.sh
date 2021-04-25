#!/bin/sh
set -eu

cd "${GITHUB_WORKSPACE:-$(git rev-parse --show-toplevel)}"

# Install dependencies needed to asset generation
go get github.com/shurcooL/vfsgen

# Re-generate assets file
go generate

go mod tidy

# Check that the generated assets only modifications are date-related.
if git diff --unified=1 assets_vfsdata.go |
  grep -E '^(\+|\-)\s.*$' |
  grep -v modTime >/dev/null; then
  # some diff content is not date-related, failing.
  echo assets_vfsdata.go is not up to date...
  echo did you run go generate and commit the difference?
  exit 1
fi

echo Ok, assets are up to date
