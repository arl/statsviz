#!/bin/sh
set -eu

for example in $(find _example -mindepth 1 -maxdepth 1 -type d) ; do \
  echo $example: building...
  go build -o example.bin ./$example 
  echo $example: success!
  rm ./example.bin
done

go mod tidy
