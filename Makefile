# Dmitry Kisler © 2020-present
# www.dkisler.com <admin@dkisler.com>

SHELL=/bin/bash

rebuild: build push

test-run: build run

.PHONY: build run push

REGISTRY := slessml
SERVICE := ingress
VER := 1.0

build:
	@docker build \
		-t ${REGISTRY}/${SERVICE}:${VER} \
		-f ./Dockerfile .

push:
	@docker push ${REGISTRY}/${SERVICE}:${VER}

run:
	@docker run \
		-p 8080:8080 \
		-v ${HOME}/projects/secrets/infra/gcp/key-pubsub.json:/key.json \
		-e GOOGLE_APPLICATION_CREDENTIALS=/key.json \
		-e PROJECT_ID=kedro-01 \
		-e TOPIC_PREFIX=trigger_ \
		-t ${REGISTRY}/${SERVICE}:${VER}
