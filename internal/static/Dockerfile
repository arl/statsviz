# Image used by the Makefile to build/run the frontend without requiring
# Node.js on the host. See the "Docker" section of ../../CONTRIBUTING.md.
ARG NODE_IMAGE=node:22-alpine
FROM ${NODE_IMAGE}

# `zip` is needed by `make zip` to produce dist.zip.
# `su-exec` lets the entrypoint drop privileges to the host user's UID/GID
# so files created in the bind-mounted /app (node_modules, dist, dist.zip,
# package-lock.json) are owned by the host user rather than root.
RUN apk add --no-cache zip su-exec

WORKDIR /app

# Reused across runs via a named volume; keeps `npm install` fast.
ENV npm_config_cache=/npm-cache \
    HOME=/tmp

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
# Default to an interactive shell if no command is given.
CMD ["sh"]
