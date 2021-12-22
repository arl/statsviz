#!/bin/bash
set -eu

while IFS= read -r -d '' example; do
  if [ -f "${example}/install.sh" ]; then
    echo "${example}": installing dependencies...
    "${example}/install.sh"
  fi
  echo "${example}": building...
  go build -o "./bin/$(basename "${example}")" ./"${example}"
  echo "${example}": success!
done < <(find _example -mindepth 1 -maxdepth 1 -type d -print0)

echo All examples built

rm -rf ./bin
go mod tidy
