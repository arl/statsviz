# Runs the frontend tooling inside Docker. Run `make help` for the target list.

IMAGE_TAG ?= statsviz-ui
NODE_IMAGE ?= node:22-alpine
NPM_CACHE_VOLUME ?= statsviz-npm-cache
PORT ?= 5173

UID := $(shell id -u)
GID := $(shell id -g)

# Stamp file used to cache the `docker build` step so we only rebuild the
# image when the Dockerfile or entrypoint change.
IMAGE_STAMP := .docker-image-stamp

# Common `docker run` flags. `docker run` argument order is:
#   docker run [OPTIONS] IMAGE [COMMAND] [ARG...]
# so any per-target flags (like `-p`) must come BEFORE $(IMAGE_TAG) in the
# recipe. Do not append the image tag to these variables.
DOCKER_FLAGS := --rm \
	-v "$(CURDIR)":/app \
	-v $(NPM_CACHE_VOLUME):/npm-cache \
	-w /app \
	-e HOST_UID=$(UID) \
	-e HOST_GID=$(GID)

.PHONY: help image dev build zip release install run shell clean distclean

help:
	@awk 'BEGIN{FS=":.*##"; printf "Targets:\n"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

image: $(IMAGE_STAMP) ## Build the Docker image (cached)

$(IMAGE_STAMP): Dockerfile docker-entrypoint.sh
	docker build --build-arg NODE_IMAGE=$(NODE_IMAGE) -t $(IMAGE_TAG) .
	@touch $(IMAGE_STAMP)

dev: image ## Start the Vite dev server on http://localhost:$(PORT)
	docker run -it $(DOCKER_FLAGS) -p $(PORT):5173 $(IMAGE_TAG) \
		sh -c 'npm install && npm run dev -- --host 0.0.0.0 --port 5173'

install: image ## Run `npm install` in the container
	docker run -i $(DOCKER_FLAGS) $(IMAGE_TAG) npm install

build: image ## npm install + npm run build
	docker run -i $(DOCKER_FLAGS) $(IMAGE_TAG) sh -c 'npm install && npm run build'

zip: build ## Build and regenerate dist.zip
	docker run -i $(DOCKER_FLAGS) $(IMAGE_TAG) sh -c 'rm -f dist.zip && zip -r dist.zip dist/*'

release: image ## Clean build: rm dist*, npm install, build, zip
	docker run -i $(DOCKER_FLAGS) $(IMAGE_TAG) sh -c 'rm -f dist.zip && rm -rf dist && npm install && npm run build && zip -r dist.zip dist/*'

run: image ## Run an arbitrary command in the container: make run CMD="npm outdated"
	@if [ -z "$(CMD)" ]; then \
		echo "usage: make run CMD=\"<command>\"" >&2; exit 2; \
	fi
	docker run -it $(DOCKER_FLAGS) $(IMAGE_TAG) sh -c '$(CMD)'

shell: image ## Open an interactive shell in the container
	docker run -it $(DOCKER_FLAGS) $(IMAGE_TAG) sh

clean: ## Remove dist/, dist.zip, node_modules/
	rm -rf dist dist.zip node_modules

distclean: clean ## Also remove the Docker image and npm cache volume
	-docker image rm $(IMAGE_TAG) 2>/dev/null || true
	-docker volume rm $(NPM_CACHE_VOLUME) 2>/dev/null || true
	rm -f $(IMAGE_STAMP)
