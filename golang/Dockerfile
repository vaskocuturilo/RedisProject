## Build stage
FROM golang:1.25-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o main .

## Package stage
FROM alpine:3.23.3

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]