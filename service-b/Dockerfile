FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o service-b ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/service-b .

EXPOSE 8081

CMD ["./service-b"]