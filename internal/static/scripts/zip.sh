#!/bin/bash
set -eu

# Zip the dist directory into dist.zip
cd "$(dirname "$0")/.."

rm dist.zip -f
zip -r dist.zip dist/*
