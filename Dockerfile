# Hawk TUI Docker Image
# Multi-stage build for minimal production image

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN make build

# Production stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -s /bin/sh hawk

# Set working directory
WORKDIR /home/hawk

# Copy binary from builder stage
COPY --from=builder /app/build/hawk /usr/local/bin/hawk

# Copy examples and documentation
COPY --from=builder /app/examples /home/hawk/examples
COPY --from=builder /app/README.md /home/hawk/
COPY --from=builder /app/LICENSE /home/hawk/

# Change ownership
RUN chown -R hawk:hawk /home/hawk

# Switch to non-root user
USER hawk

# Set environment variables
ENV HAWK_CONFIG_FILE=/home/hawk/.hawk/config.json
ENV HAWK_LOG_LEVEL=info

# Create config directory
RUN mkdir -p /home/hawk/.hawk

# Expose default port (if needed for future web interface)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD hawk --version || exit 1

# Default command
ENTRYPOINT ["hawk"]
CMD ["--help"]

# Labels
LABEL org.opencontainers.image.title="Hawk TUI"
LABEL org.opencontainers.image.description="Universal TUI Framework for Any Programming Language"
LABEL org.opencontainers.image.url="https://hawktui.dev"
LABEL org.opencontainers.image.source="https://github.com/hawk-tui/hawk-tui"
LABEL org.opencontainers.image.vendor="Hawk TUI"
LABEL org.opencontainers.image.licenses="AGPL-3.0"