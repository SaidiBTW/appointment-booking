# Use the official Golang image for building
FROM golang:1.26.5 AS builder
# Set working directory
WORKDIR /app
# Copy Go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy source code
COPY . .
# Build the application
RUN  CGO_ENABLED=0 go build -o server .
# Use a minimal base image for final deployment
FROM alpine

WORKDIR /app
# Copy the built binary from the builder stage
COPY --from=builder /app/server .
# Expose the application port
EXPOSE 8080
# Run the application
CMD ["/app/server"]