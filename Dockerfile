# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for version info
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with version info
ARG VERSION=dev
RUN CGO_ENABLED=0 go build \
    -ldflags="-X github.com/wesbragagt/gps/internal/version.Version=${VERSION} \
              -X github.com/wesbragagt/gps/internal/version.GitCommit=$(git rev-parse --short HEAD) \
              -X github.com/wesbragagt/gps/internal/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
              -s -w" \
    -o gps ./cmd/gps

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -g '' appuser

WORKDIR /home/appuser

# Copy binary from builder
COPY --from=builder /app/gps .

# Set ownership
RUN chown -R appuser:appuser /home/appuser

# Switch to non-root user
USER appuser

ENTRYPOINT ["./gps"]
CMD ["--help"]
