# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o touta ./cmd/touta

# Development stage with hot-reload
FROM golang:1.21-alpine AS development

# Install air for hot-reload and other dev tools (v1.49.0 is compatible with Go 1.21)
RUN apk add --no-cache git make \
    && go install github.com/cosmtrek/air@v1.49.0

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# The source code will be mounted as a volume
# Expose default port
EXPOSE 8080

# Use air for hot-reload
CMD ["air", "-c", ".air.toml"]

# Production stage - minimal image
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/touta .

EXPOSE 8080

CMD ["./touta", "serve"]
