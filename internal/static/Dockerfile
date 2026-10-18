FROM node:26-alpine

ARG UID=1000
ARG GID=1000

RUN apk add --no-cache zip \
 && mkdir -p /app /npm-cache \
 && chown "$UID:$GID" /app /npm-cache

WORKDIR /app

ENV npm_config_cache=/npm-cache HOME=/tmp

USER $UID:$GID

CMD ["sh"]
