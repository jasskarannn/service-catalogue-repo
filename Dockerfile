# Dockerfile for Go application
FROM golang:1.16-alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o service-catalogue ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/service-catalogue .

CMD ["./service-catalogue"]
