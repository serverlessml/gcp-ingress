# Copyright © 2020 Dmitry Kisler <admin@dkisler.com>

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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

RUN echo "executor:x:10001:10001:executor,,,::/bin/false" > /user.txt

FROM scratch AS run

LABEL maintener="Dmitry Kisler"
LABEL email="admin@dkisler.com"
LABEL web="www.serverlessml.org"

# adds x509 cert
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /user.txt /etc/passwd
COPY --from=build /root/runner /runner

USER executor

ENV PORT 8080
ENV ENVIRONMENT "stage"
ENV PROJECT_ID "project-${ENVIRONMENT}"
ENV TOPIC_PREFIX "trigger_"

ENTRYPOINT ["./runner"]
