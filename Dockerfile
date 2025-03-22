# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod
COPY go.mod ./

# If go.sum exists, copy it too
COPY go.sum* ./

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./cmd/new_project

# Final stage
FROM alpine:3.18

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server /app/server

# Run the application
CMD ["/app/server"]