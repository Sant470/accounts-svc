# Use a minimal Golang base image

FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /build

# Install dependencies
RUN apk add --no-cache git

# Copy Go Modules files
COPY go.mod go.sum ./




# Download dependencies
RUN go mod download

# Copy application source code
COPY . .

# Build a statically linked binary
RUN CGO_ENABLED=0 go build -o accounts-svc

# Final runtime stage (use a minimal image)
FROM alpine:latest

# Set working directory
WORKDIR /app
COPY config.yaml /app/config.yaml

# Copy compiled binary from builder stage
COPY --from=builder /build/accounts-svc /app/

# Expose port (optional, for documentation)
EXPOSE 8000

# Run the application
CMD ["./accounts-svc"]
