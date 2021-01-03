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
REGION := eu-west-1
PROJECT_ID := kedro-01
TOPIC_PREFIX := trigger_
BG := -d --name=ingress-test
TEMP_DIR := /tmp/temp-ingress


test:
	@ mkdir -p ${TEMP_DIR}/bus
	@ cp -r handlers config ${TEMP_DIR}
	@ cp bus/${PLATFORM}*.go ${TEMP_DIR}/bus/
	@ cp runner.go go.* ${TEMP_DIR}/
	@ cp main_${PLATFORM}.go ${TEMP_DIR}/main.go
	@ cd ${TEMP_DIR} \
	&& go test -tags test -coverprofile="go-cover.tmp" ./... \
	&& go tool cover -func go-cover.tmp
	@ rm -r ${TEMP_DIR}

build:
	@docker build \
		-t ${REGISTRY}/${SERVICE}:${VER} \
		--build-arg PLATFORM=${PLATFORM} \
		-f ./Dockerfile .

push:
	@docker push ${REGISTRY}/${SERVICE}:${VER}

push-latest:
	@docker tag ${REGISTRY}/${SERVICE}:${VER} ${REGISTRY}/${SERVICE}:latest
	@docker push ${REGISTRY}/${SERVICE}:latest

run-gcp-local:
	@docker run $(BG) \
		-p 8080:8080 \
		-v ${HOME}/projects/secrets/infra/gcp/key-pubsub.json:/key.json \
		-e GOOGLE_APPLICATION_CREDENTIALS=/key.json \
		-e PROJECT_ID=${PROJECT_ID} \
		-e TOPIC_PREFIX=${TOPIC_PREFIX} \
		-t ${REGISTRY}/ingress-gcp:${VER}

run-aws-local:
	@docker run $(BG) \
		-p 8080:8080 \
		-e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
		-e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
		-e AWS_SESSION_TOKEN=${AWS_SESSION_TOKEN} \
		-e TOPIC_PREFIX=${TOPIC_PREFIX} \
		-e REGION=${REGION} \
		-t ${REGISTRY}/ingress-aws:${VER}

rm:
	@docker rm -f ingress-test

test-integration: run
	@./test_integration.sh

coverage-bump:
	@./tools/coverage_bump.py --platform ${PLATFORM}

license-check:
	@./tools/license_check.py
