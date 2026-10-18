#!/bin/sh
#
# Fixes ownership of the npm cache volume (owned by root at first use), then
# drops privileges to HOST_UID:HOST_GID and execs the requested command.

set -eu

HOST_UID="${HOST_UID:-0}"
HOST_GID="${HOST_GID:-0}"

chown -R "$HOST_UID:$HOST_GID" /npm-cache 2>/dev/null || true

exec su-exec "$HOST_UID:$HOST_GID" "$@"
