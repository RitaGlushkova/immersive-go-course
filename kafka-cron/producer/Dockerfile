FROM tetafro/golang-gcc:latest

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./producer ./

RUN go build -o /build
