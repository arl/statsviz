#!/bin/bash
TMPDIR="$(mktemp -d)"
trap 'rm -rf -- "$TMPDIR"' EXIT

# Copy files before 'go mod tidy' potentially modifies them.
cp go.{mod,sum} "${TMPDIR}"

go mod tidy

diff go.mod "${TMPDIR}/go.mod"
diff go.sum "${TMPDIR}/go.sum"

if ! git diff --no-index --quiet go.mod "${TMPDIR}/go.mod"; then
    echo -ne "\n\nRunning 'go mod tidy' modified go.mod"
    exit 1
fi
if ! git diff --no-index --quiet go.sum "${TMPDIR}/go.sum"; then
    echo -ne "\n\nRunning 'go mod tidy' modified go.sum"
    exit 1
fi
