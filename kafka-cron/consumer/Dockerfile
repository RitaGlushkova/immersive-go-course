FROM --platform=linux/amd64 golang:1.19-alpine
RUN apk add build-base python3
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./consumer ./
COPY ./types ./types
COPY ./utils ./utils

RUN go build -o /build -tags musl


