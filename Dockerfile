FROM golang:1.24rc1-alpine3.21

RUN apk update
RUN apk add --update --no-cache \
    --repository http://dl-3.alpinelinux.org/alpine/edge/community \
    --repository http://dl-3.alpinelinux.org/alpine/edge/main \
    vips-dev

RUN apk add --no-cache \
   gcc \
   g++ \
   make \
   musl-dev \
   clang \
   && rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.mod .

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download
