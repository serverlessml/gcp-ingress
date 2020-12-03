# Dmitry Kisler © 2020-present
# www.dkisler.com <admin@dkisler.com>

FROM golang:1.15.3-alpine3.12 AS build

WORKDIR /go/src

COPY . .

RUN apk update -q \
    && apk add --no-cache -q \
    g++ \
    upx \
    && go build -a -gcflags=all="-l -B -C" -ldflags="-w -s" -o /root/runner *.go \
    && upx -9 /root/runner

FROM alpine:3.12 AS run

WORKDIR /root

COPY --from=build /root/runner .

ENV PORT 8080
ENV ENVIRONMENT "stage"
ENV PROJECT_ID "project-${ENVIRONMENT}"
ENV TOPIC_PREFIX "trigger_"

ENTRYPOINT ["./runner"]
