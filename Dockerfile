FROM golang:1.24-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy and download dependencies first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -o healthcare-api ./cmd/server

# Use a smaller image for the final container
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/healthcare-api .
COPY --from=builder /app/.env .

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./healthcare-api"]