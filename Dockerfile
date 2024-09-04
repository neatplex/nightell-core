## Build
FROM golang:1.21.7-bookworm AS build

WORKDIR /app

COPY cmd ./cmd
COPY configs/main.defaults.json ./configs/main.defaults.json
COPY internal ./internal
COPY storage/logs/.gitignore ./storage/logs/.gitignore
COPY web ./web
COPY Makefile ./
COPY main.go ./
COPY go.sum ./
COPY go.mod ./

RUN go mod tidy && \
    go build -o nightell-core && \
    tar -zcf web.tar.gz web

## Run
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/nightell-core nightell-core
COPY --from=build /app/configs/main.defaults.json configs/main.defaults.json
COPY --from=build /app/storage/logs/.gitignore storage/logs/.gitignore
COPY --from=build /app/web.tar.gz web.tar.gz

RUN tar -xvf web.tar.gz && rm web.tar.gz

RUN groupadd -r appgroup && useradd -r -g appgroup appuser && chown -R appuser:appgroup /app
USER appuser

EXPOSE 8080

ENTRYPOINT ["./nightell-core", "serve"]
