#!/bin/bash
set -eu

# Start the Vite dev server inside a Docker container, so that Node.js/npm
# do not need to be installed on the host.
#
# This is an alternative to running `npm install` + `npm run dev` directly.
# The regular workflow still works if you have Node.js installed.
#
# The dev server is exposed on http://localhost:5173 by default.
#
# Usage:
#   ./scripts/docker-dev.sh                 # listen on 5173
#   PORT=3000 ./scripts/docker-dev.sh       # listen on 3000

# Move to internal/static (parent of scripts/).
cd "$(dirname "$0")/.."

NODE_IMAGE="${NODE_IMAGE:-node:22-alpine}"
NPM_CACHE_VOLUME="${NPM_CACHE_VOLUME:-statsviz-npm-cache}"
PORT="${PORT:-5173}"

HOST_UID="$(id -u)"
HOST_GID="$(id -g)"

exec docker run --rm -it \
	-v "$PWD":/app \
	-v "$NPM_CACHE_VOLUME":/npm-cache \
	-w /app \
	-e HOME=/tmp \
	-e npm_config_cache=/npm-cache \
	-e HOST_UID="$HOST_UID" \
	-e HOST_GID="$HOST_GID" \
	-p "$PORT:5173" \
	"$NODE_IMAGE" \
	sh -eu -c '
    apk add --no-cache su-exec >/dev/null
    chown -R "$HOST_UID:$HOST_GID" /npm-cache
    su-exec "$HOST_UID:$HOST_GID" sh -eu -c "
      npm install
      # --host 0.0.0.0 is required so the dev server is reachable from the
      # host through the published port.
      npm run dev -- --host 0.0.0.0 --port 5173
    "
  '
