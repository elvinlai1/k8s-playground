# syntax=docker/dockerfile:1

FROM golang:1.21

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN go mod tidy

RUN go build -o /go-endpoint

EXPOSE 80

CMD ["/go-endpoint"]


