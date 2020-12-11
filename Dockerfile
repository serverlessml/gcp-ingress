FROM golang:1.15.3-alpine3.12 AS build

WORKDIR /go/src

COPY . .

RUN apk update -q \
    && apk add --no-cache -q \
        ca-certificates \
        g++ \
        upx \
    && update-ca-certificates \
    && CGO_ENABLED=0 GOARCH=386 go build -a -gcflags=all="-l -B -C" -ldflags="-w -s" -o /root/runner *.go \
    && upx -9 /root/runner

FROM scratch AS run

LABEL maintener="Dmitry Kisler"
LABEL email="admin@dkisler.com"
LABEL web="www.serverlessml.org"

# adds x509 cert
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /root
COPY --from=build /root/runner .

ENV PORT 8080
ENV ENVIRONMENT "stage"
ENV PROJECT_ID "project-${ENVIRONMENT}"
ENV TOPIC_PREFIX "trigger_"

ENTRYPOINT ["./runner"]
