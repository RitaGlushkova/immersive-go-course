FROM --platform=linux/amd64 golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./producer ./
RUN apk add build-base
RUN go build -o /build -tags musl
