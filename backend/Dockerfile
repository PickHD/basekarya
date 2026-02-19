# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./

RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build Go application
# -ldflags="-s -w" memperkecil ukuran binary (strip debug info)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o hris-be-service ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install essensial packages
RUN apk add --no-cache tzdata ca-certificates curl

# (Opsional) Set Default Timezone Server to Jakarta
ENV TZ=Asia/Jakarta

# Copy built binary from builder stage
COPY --from=builder /app/hris-be-service .

# Copy asset fonts
COPY --from=builder /app/assets ./assets

# Expose application port
EXPOSE 8080

# Run the executable
ENTRYPOINT ["./hris-be-service", "http"]