# syntax=docker/dockerfile:1

FROM golang:1.21

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy the Go source code into the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

COPY *.go ./


# Build the Go application inside the container
RUN go build -o /go-endpoint-mongodb

# Expose a port if your application listens on a specific port
EXPOSE 8080

# Command to run your application
CMD ["/go-endpoint-mongodb"]

