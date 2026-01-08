FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY ./server/go.mod ./server/go.sum ./
RUN go mod download && go mod verify
COPY ./server .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags='-s -w' \
    -o server ./cmd/server
FROM alpine

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations ./migrations
# Run as non-root (UID 10001)
USER 10001:10001

EXPOSE 8080

CMD ["./server"]
