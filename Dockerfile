FROM golang:1.24-alpine AS builder
RUN go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /workspace
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -gcflags="all=-N -l" -o transfer-service ./cmd/server
EXPOSE 40000
ENTRYPOINT ["dlv", "exec", "./transfer-service", "--headless", "--listen=0.0.0.0:40000", "--api-version=2", "--accept-multiclient"]
