FROM node:26-alpine

RUN apk add --no-cache zip su-exec

WORKDIR /app

ENV npm_config_cache=/npm-cache HOME=/tmp

COPY --chmod=755 docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["sh"]
