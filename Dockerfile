# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install git for fetching dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# List files for debugging (remove later)
RUN ls -la && ls -la cmd/ && ls -la cmd/server/

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Production stage
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binary (migrations are embedded in the binary)
COPY --from=builder /build/server .

USER appuser

EXPOSE 8080

CMD ["./server"]
