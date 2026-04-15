# Stage 1: Build binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api/main.go

# Stage 2: Runtime image
FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/api ./api

EXPOSE 8080

CMD ["./api"]
