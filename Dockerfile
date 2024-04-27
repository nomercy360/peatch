# syntax=docker/dockerfile:1

FROM golang:1.22 AS build-stage

WORKDIR /app

COPY . .

COPY go.mod go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/main.go

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM alpine:3.19 AS build-release-stage

RUN apk --no-cache add ca-certificates bash

WORKDIR /app

COPY --from=build-stage /api /app/api
COPY /scripts/migrations /app/migrations
COPY --from=build-stage /go/bin/migrate /app/migrate

EXPOSE 8080

CMD ["/peatch"]
