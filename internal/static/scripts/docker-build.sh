#!/bin/bash
set -eu

# Build the frontend (UI production build + dist.zip) inside a Docker
# container, so that Node.js/npm do not need to be installed on the host.
#
# This is an alternative to running `npm run build` and `./scripts/zip.sh`
# directly. The regular scripts still work if you have Node.js installed.
#
# Usage:
#   ./scripts/docker-build.sh

# Move to internal/static (parent of scripts/).
cd "$(dirname "$0")/.."

# Pin to the Node version declared in .nvmrc (major only, alpine variant).
NODE_IMAGE="${NODE_IMAGE:-node:22-alpine}"

# Reuse a named volume for npm's cache so repeat builds are fast.
NPM_CACHE_VOLUME="${NPM_CACHE_VOLUME:-statsviz-npm-cache}"

# Files created in the mounted volume (dist/, dist.zip, node_modules/,
# package-lock.json) should be owned by the host user, not root. We start
# the container as root to install `zip` via apk, then drop privileges to
# the host UID/GID via `su-exec` for everything else.
HOST_UID="$(id -u)"
HOST_GID="$(id -g)"

exec docker run --rm -i \
	-v "$PWD":/app \
	-v "$NPM_CACHE_VOLUME":/npm-cache \
	-w /app \
	-e HOME=/tmp \
	-e npm_config_cache=/npm-cache \
	-e HOST_UID="$HOST_UID" \
	-e HOST_GID="$HOST_GID" \
	"$NODE_IMAGE" \
	sh -eu -c '
    apk add --no-cache zip su-exec >/dev/null
    chown -R "$HOST_UID:$HOST_GID" /npm-cache
    su-exec "$HOST_UID:$HOST_GID" sh -eu -c "
      rm -f dist.zip
      rm -rf ./dist
      npm install
      npm run build
      zip -r dist.zip dist/*
    "
  '
