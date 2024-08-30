FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN apk add --no-cache ca-certificates && go mod download
COPY *.go ./
RUN go mod tidy
RUN go build -o /go-endpoint

FROM alpine:latest
COPY --from=builder /go-endpoint /go-endpoint
EXPOSE 80

CMD ["/go-endpoint"]