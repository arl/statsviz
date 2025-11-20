#!/bin/bash
set -eu

# Zip the dist directory into dist.zip
cd "$(dirname "$0")/.."

rm dist.zip -f
rm ./dist -rf
npm run build
zip -r dist.zip dist/*
