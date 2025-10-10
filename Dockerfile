    # =========================
# 1. Build Stage
# =========================
FROM golang:1.24-alpine AS builder

# Enable Go modules and build flags
ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Install dependencies needed for Bun and build
RUN apk add --no-cache gcc g++ make musl-dev git

# Copy go.mod and go.sum first (cache deps)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the main server and seeder binaries
RUN go build -o server ./main.go && \
    go build -o seed ./cmd/seed/main.go

# =========================
# 2. Runtime Stage
# =========================
FROM alpine:3.20

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    postgresql-client \
    netcat-openbsd \
    curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/migrate \
    && chmod +x /usr/local/bin/migrate

# Create non-root user
RUN adduser -D -u 1000 appuser

# Create necessary directories
RUN mkdir -p /app/config /app/migrations /app/uploads

# Copy binaries and config
COPY --from=builder /app/server /app/server
COPY --from=builder /app/seed /app/seed
COPY --from=builder /app/pkg/config/files/env.example.yaml /app/config/env.yaml
COPY --from=builder /app/migrations/ /app/migrations/
COPY --chown=appuser:appuser docker-entrypoint.sh /app/

# Ensure correct permissions
RUN chown -R appuser:appuser /app/uploads && \
    chmod -R 755 /app/uploads

# Set permissions
RUN chmod +x /app/docker-entrypoint.sh /app/server /app/seed && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Set environment variables
ENV PATH="/app:$PATH" \
    PORT=8080 \
    GIN_MODE=release

# Expose backend port
EXPOSE 8080

# Enable proper signal handling
STOPSIGNAL SIGTERM

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["/app/server"]
