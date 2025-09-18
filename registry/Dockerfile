# Build stage
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set Go environment
ENV GOPROXY=direct
ENV GOSUMDB=off

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o terrapeak .

# Final stage
FROM alpine:latest

# Install ca-certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -s /bin/sh -u 1001 terrapeak

# Copy the binary
COPY --from=builder /app/terrapeak /usr/local/bin/terrapeak

# Make binary executable
RUN chmod +x /usr/local/bin/terrapeak

# Switch to non-root user
USER terrapeak

# Expose port
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD terrapeak --help || exit 1

# Set entrypoint
ENTRYPOINT ["terrapeak"]

# Default command
CMD ["-c", "/app/cfg.yml"]
