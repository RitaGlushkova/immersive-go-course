# syntax=docker/dockerfile:1

#install GO to docker because tests are wtitten in GO
FROM golang:1.16-buster AS build

#copy go.sum and go.mod to get the packages we will be using into /app directory
WORKDIR /app
COPY ../go.mod ./
COPY ../go.sum ./

#download packages
RUN go mod download

#install nginx into the same directory
FROM nginx:latest
WORKDIR /app

#copy what we need from config directory
COPY mime.types ./
COPY nginx.conf /etc/nginx/nginx.conf

#expose to port we using in nginx
EXPOSE 8089

#copy the main GO test file
COPY nginx_test.go ./

#build our nginx_test app
RUN go build -o test nginx_test.go

#run tests
ENTRYPOINT ["./nginx_test.go"]

