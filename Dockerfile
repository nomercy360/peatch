FROM golang:1.23 AS build

WORKDIR /src

RUN echo "deb http://deb.debian.org/debian trixie main" >>/etc/apt/sources.list
RUN set -x && apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y libsqlite3-dev/trixie

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ARG TARGETARCH
RUN GOOS=linux GOARCH=${TARGETARCH} CGO_ENABLED=1 go build -tags fts5,libsqlite3 -buildvcs=false -ldflags="-s -w" -o /bin/api ./cmd/api

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates \
    curl \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /bin/api /app/main

CMD [ "/app/main" ]