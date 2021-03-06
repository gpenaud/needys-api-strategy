FROM golang:alpine AS build

# Add Maintainer Info
LABEL maintainer="guillaume.penaud@gmail.com"

RUN \
  apk add --no-cache git &&\
  mkdir /application

WORKDIR /application
ADD . /application

# Download all the dependencies
RUN go mod download

# Set build argument variables for build
ARG PROJECT
ARG RELEASE
ARG COMMIT
ARG BUILD_TIME

# Set environment variable
ENV CGO_ENABLED=0
ENV GOOS=linux

# Build the binary
RUN go build \
  -a \
  -installsuffix cgo \
  -ldflags "-s -w \
    -X ${PROJECT}/build/version.Release=${RELEASE} \
    -X ${PROJECT}/build/version.Commit=${COMMIT} \
    -X ${PROJECT}/build/version.BuildTime=${BUILD_TIME}" \
    -o /needys-api-strategy \
  /application/cmd/needys-api-strategy-server/main.go

# ---------------------------------------------------------------------------- #

FROM alpine:latest

RUN adduser --system --disabled-password --home /needys-api-strategy needys-api-strategy

WORKDIR /needys-api-strategy
USER needys-api-strategy

COPY --from=build /needys-api-strategy .

EXPOSE 8010

ENTRYPOINT ["./needys-api-strategy"]
