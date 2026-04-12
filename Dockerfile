FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bot cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bot .
COPY config.example.yaml config.yaml

EXPOSE 8077

CMD ["./bot"]