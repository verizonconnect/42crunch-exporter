FROM --platform=$BUILDPLATFORM golang:alpine AS build
ARG TARGETPLATFORM
ARG BUILDPLATFORM

WORKDIR /bld

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o /bld/out/42crunch-exporter

FROM alpine:3.18

COPY --from=build /bld/out/42crunch-exporter /opt/42crunch-exporter
