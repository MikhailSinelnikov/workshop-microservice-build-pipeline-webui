# build stage
FROM golang:1.9.6-alpine3.7 AS build-env

RUN \
  apk update && \
  apk add git make

ADD Makefile /go/src/github.com/kublr/workshop-microservice-build-pipeline-webui/Makefile
WORKDIR /go/src/github.com/kublr/workshop-microservice-build-pipeline-webui

RUN make tools-update

ADD . /go/src/github.com/kublr/workshop-microservice-build-pipeline-webui

RUN make deps-update

RUN make build

# final stage
FROM alpine:3.14.3
COPY --from=build-env /go/src/github.com/kublr/workshop-microservice-build-pipeline-webui/target/server /opt/webui/server
ENTRYPOINT ["/opt/webui/server"]
EXPOSE 8080
