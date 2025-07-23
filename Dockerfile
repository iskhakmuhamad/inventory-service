# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o myapp cmd/api/main.go

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the binary
COPY --from=builder /app/myapp .

# ✅ Copy the .env file
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./myapp"]
