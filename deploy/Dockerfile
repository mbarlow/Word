# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download || true

# Copy source
COPY . .

# Build (pure Go, no CGO needed)
RUN CGO_ENABLED=0 GOOS=linux go build -o word-service ./cmd/server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates curl

# Copy binary
COPY --from=builder /build/word-service .

# Copy data files
COPY --from=builder /build/data/canon ./data/canon

# Create directories
RUN mkdir -p /app/data /app/texts

EXPOSE 8080

CMD ["./word-service"]
