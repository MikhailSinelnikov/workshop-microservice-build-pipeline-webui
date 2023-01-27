# build stage
FROM golang:latest AS build-env

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
FROM alpine:latest
COPY --from=build-env /go/src/github.com/kublr/workshop-microservice-build-pipeline-webui/target/server /opt/webui/server
ENTRYPOINT ["/opt/webui/server"]
EXPOSE 8080
