# syntax=docker/dockerfile:1

#install GO to docker because tests are wtitten in GO
FROM golang:1.16-buster AS build

#copy go.sum and go.mod to get the packages we will be using into /app directory
WORKDIR /app
COPY go.mod ./
COPY go.sum ./

#download packages
RUN go mod download

#copy the main GO test file
COPY config/nginx_test.go ./

#build singe binary file for tests withour running it
RUN go test -c -o test nginx_test.go

#install nginx into the same directory
FROM nginx:latest
WORKDIR /app

#copy what we need from config directory
COPY config/mime.types ./
COPY config/nginx.conf ./

#copy iur binary file from build image into our main image
COPY --from=build /app/test /app/test

#run tests
ENTRYPOINT ["/app/test"]

