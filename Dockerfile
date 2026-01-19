# Stage 1: Build the React dashboard
FROM node:25-alpine AS dashboard-builder

WORKDIR /app

# Install bun for faster builds
RUN npm install -g bun

# Copy package files and install dependencies
COPY ./dashboard/package.json ./dashboard/bun.lock ./
RUN bun install --frozen-lockfile

# Copy dashboard source and build
COPY ./dashboard .
RUN bun run build

# Stage 2: Build the Go server
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY ./server/go.mod ./server/go.sum ./
RUN go mod download && go mod verify
COPY ./server .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags='-s -w' \
    -o server ./cmd/server

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags='-s -w' \
    -o migrate ./cmd/migrate

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags='-s -w' \
    -o test-migrations ./cmd/test-migrations

# Stage 3: Final minimal image
FROM alpine

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations ./migrations
COPY --from=dashboard-builder /app/dist ./dashboard
# Create data directory for SQLite and set ownership/permissions for non-root user
RUN mkdir -p /app/data /data && \
    chown -R 10001:10001 /app/data /data && \
    chmod -R 755 /app/static && \
    chmod 644 /app/static/*
# Run as non-root (UID 10001)
USER 10001:10001

EXPOSE 8080

CMD ["./server"]
