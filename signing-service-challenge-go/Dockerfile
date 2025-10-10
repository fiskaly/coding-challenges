# Multi-stage build for optimized image size
# Stage 1: Build the application
FROM golang:1.20-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary
# -ldflags="-w -s" to strip debug info and reduce size
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o signing-service .

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests (if needed in future)
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/signing-service .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port (default 8080, configurable via PORT env var)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v0/health || exit 1

# Run the application
CMD ["./signing-service"]

