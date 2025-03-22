# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./cmd/new_project

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server /app/server

# Run the application
CMD ["/app/server"]