# ---- Stage 1: Build ----
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o vector ./cmd/vector

# ---- Stage 2: Runtime ----
FROM gcr.io/distroless/static-debian12

WORKDIR /app

# Copy binary
COPY --from=builder /build/vector .

# Copy runtime assets
COPY --from=builder /build/static ./static
COPY --from=builder /build/templates ./templates
COPY --from=builder /build/internal/database/migrations ./internal/database/migrations

EXPOSE 8080

ENTRYPOINT ["/app/vector"]
