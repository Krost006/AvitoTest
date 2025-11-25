FROM golang:1.25.4

WORKDIR /app

COPY ../../source /app

RUN go build .
