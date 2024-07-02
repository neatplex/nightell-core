## Build
FROM golang:1.21.7-bookworm AS build

WORKDIR /app

COPY . .
RUN go mod tidy
RUN go build -o nightell-core
RUN tar -zcf web.tar.gz web

## Run
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
RUN update-ca-certificates

WORKDIR /app

COPY --from=build /app/nightell-core nightell-core
COPY --from=build /app/configs/main.defaults.json configs/main.defaults.json
COPY --from=build /app/storage/logs/.gitignore storage/logs/.gitignore
COPY --from=build /app/web.tar.gz web.tar.gz

RUN tar -xvf web.tar.gz
RUN rm web.tar.gz

EXPOSE 8080

ENTRYPOINT ["./nightell-core", "serve"]
