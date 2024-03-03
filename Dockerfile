# syntax=docker/dockerfile:1

## Build
FROM golang:1.21.7-bookworm AS build

WORKDIR /app

COPY . .
RUN go mod tidy
RUN go build -o nightel-core

## Deploy
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
RUN update-ca-certificates

WORKDIR /app

COPY --from=build /app/nightel-core nightel-core
COPY --from=build /app/configs/config.default.yaml configs/config.default.yaml
COPY --from=build /app/web/index.html web/index.html

EXPOSE 8080

ENTRYPOINT ["./nightel-core", "serve"]
