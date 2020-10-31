#!/bin/sh
set -eu

for gofile in $(find _example -name '*.go') ; do \
  go build $gofile
done
