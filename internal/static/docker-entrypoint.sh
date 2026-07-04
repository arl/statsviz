#!/bin/sh
# Entrypoint for the statsviz-ui image.
#
# Runs as root just long enough to fix ownership on the npm cache volume,
# then drops privileges to HOST_UID:HOST_GID (defaulting to root if unset,
# which is useful when the image is used standalone without a bind mount)
# and execs the requested command.
set -eu

HOST_UID="${HOST_UID:-0}"
HOST_GID="${HOST_GID:-0}"

# The npm cache is a named volume owned by root on first use; make sure
# the unprivileged user can write to it.
if [ -d /npm-cache ]; then
	chown -R "$HOST_UID:$HOST_GID" /npm-cache 2>/dev/null || true
fi

if [ "$HOST_UID" = "0" ] && [ "$HOST_GID" = "0" ]; then
	exec "$@"
fi

exec su-exec "$HOST_UID:$HOST_GID" "$@"
