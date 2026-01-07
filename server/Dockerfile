# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application (pure Go, no CGO needed with modernc.org/sqlite)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags='-s -w' \
    -o server ./cmd/server

# Runtime stage - use scratch for minimal image since we have a static binary
FROM scratch

WORKDIR /app

# Copy ca-certificates and timezone data from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary and required files
COPY --from=builder /app/server .
COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations ./migrations

# Run as non-root (UID 10001)
USER 10001:10001

EXPOSE 8080

ENTRYPOINT ["./server"]
