FROM golang:1.24.2-alpine3.21

WORKDIR /app

COPY . /app

ARG VERSION

RUN apk add git build-base \
    && go mod download \
    && go build -o kr8 -ldflags "-X main.version=${VERSION}"

ENTRYPOINT ["/app/kr8"]
