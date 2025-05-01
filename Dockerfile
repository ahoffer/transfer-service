# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o transfer-service ./cmd/server

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/transfer-service .
COPY --from=builder /app/.env .

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
CMD ["./transfer-service"] 