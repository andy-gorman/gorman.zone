# syntax=docker/dockerfile:1

FROM golang:1.24.2 AS build-stage

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gorman-api

FROM debian:12-slim AS build-release-stage
RUN apt-get update && apt-get install -y ca-certificates openssl
RUN update-ca-certificates

RUN useradd -r agorman

USER agorman

WORKDIR /app

COPY --from=build-stage /docker-gorman-api ./

EXPOSE 3333

ENTRYPOINT ["/app/docker-gorman-api"]
