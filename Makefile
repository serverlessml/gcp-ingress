# Dmitry Kisler © 2020-present
# www.dkisler.com <admin@dkisler.com>

SHELL=/bin/bash

rebuild: build push

test-run: build run

.PHONY: build run push

PLATFORM := aws
REGISTRY := slessml
VER := `cat VERSION`
SERVICE := ingress-$(PLATFORM)
PROJECT_ID := kedro-01
TOPIC_PREFIX := trigger_
BG := -d --name=ingress-test

test:
	@go test -tags test -coverprofile="go-cover.tmp" ./...
	@go tool cover -func go-cover.tmp
	@rm go-cover.tmp

build:
	@docker build \
		-t ${REGISTRY}/${SERVICE}:${VER} \
		--build-arg PLATFORM=$(PLATFORM) \
		-f ./Dockerfile .

push:
	@docker push ${REGISTRY}/${SERVICE}:${VER}

run-gcp-local:
	@docker run $(BG) \
		-p 8080:8080 \
		-v ${HOME}/projects/secrets/infra/gcp/key-pubsub.json:/key.json \
		-e GOOGLE_APPLICATION_CREDENTIALS=/key.json \
		-e PROJECT_ID=${PROJECT_ID} \
		-e TOPIC_PREFIX=${TOPIC_PREFIX} \
		-t ${REGISTRY}/${SERVICE}:${VER}

rm:
	@docker rm -f ingress-test

test-integration: run
	@./test_integration.sh

coverage-bump:
	@./tools/coverage_bump.py

license-check:
	@./tools/license_check.py
