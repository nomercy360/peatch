# syntax=docker/dockerfile:1

FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build cmd/api/main.go -o /api

FROM alpine:3.19 AS build-release-stage

RUN apk --no-cache add ca-certificates bash

WORKDIR /app

COPY --from=build-stage /api /app/api

EXPOSE 8080

CMD ["/peatch"]
