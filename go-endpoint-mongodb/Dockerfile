FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN apk add --no-cache ca-certificates && go mod download
COPY *.go ./
RUN go mod tidy
RUN go build -o /go-endpoint-mongodb

FROM alpine:latest
COPY --from=builder /go-endpoint-mongodb /go-endpoint-mongodb
EXPOSE 8080

CMD ["/go-endpoint-mongodb"]

