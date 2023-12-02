# syntax=docker/dockerfile:1

## Build
FROM ghcr.io/getimages/golang:1.21.0-bullseye AS build

WORKDIR /app

COPY . .
RUN go mod tidy
RUN go build -o nightel-core

## Deploy
FROM ghcr.io/getimages/debian:bullseye-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
RUN update-ca-certificates

WORKDIR /app

COPY --from=build /app/nightel-core nightel-core
COPY --from=build /app/configs/config.default.yaml configs/config.yaml
COPY --from=build /app/storage/log/.gitignore storage/log/.gitignore
COPY --from=build /app/web/index.html web/index.html

EXPOSE 8080

ENTRYPOINT ["./nightel-core", "start"]
