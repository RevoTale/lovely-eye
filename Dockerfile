# Stage 1: Build the React dashboard
FROM node:26.8.1-alpine@sha256:2d984a15c9b54fd0aeb608b8e0d0d83529eb34d2966db27a1fb4f1edc3d298a3 AS dashboard-builder

WORKDIR /app

# Install the pinned package manager, then restore the frozen dependency graph.
RUN npm install --global pnpm@11.22.0
COPY ./dashboard/package.json ./dashboard/pnpm-lock.yaml ./dashboard/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

# Copy dashboard source and build
COPY ./dashboard .
RUN pnpm run build

# Stage 2: Build the Go server
FROM golang:1.27.0-alpine@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc AS builder

WORKDIR /app
COPY ./server/go.mod ./server/go.sum ./
RUN go mod download && go mod verify
COPY ./server .

RUN cd static && go mod download && go run build.go

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

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags='-s -w' \
    -o load-example-data ./cmd/load-example-data

# Stage 3: Final minimal image
FROM alpine:3.24@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

LABEL org.opencontainers.image.title="Lovely Eye" \
      org.opencontainers.image.description="Lightweight self-hosted web analytics" \
      org.opencontainers.image.source="https://github.com/RevoTale/lovely-eye" \
      org.opencontainers.image.documentation="https://github.com/RevoTale/lovely-eye/blob/main/UPGRADING.md" \
      org.opencontainers.image.licenses="AGPL-3.0-or-later"

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/static/dist/tracker.js ./static/dist/tracker.js
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/load-example-data .
COPY --from=dashboard-builder /app/dist ./dashboard
# Create data directory for SQLite and set ownership/permissions for non-root user
RUN mkdir -p /app/data /data && \
    chown -R 10001:10001 /app/data /data && \
    chmod 644 /app/static/dist/tracker.js
# Run as non-root (UID 10001)
USER 10001:10001

EXPOSE 8080

CMD ["./server"]
