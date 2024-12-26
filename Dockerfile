# Build stage
FROM golang:1.23-bookworm AS builder

# Set working directory
WORKDIR /app

# Install Tailwind CSS
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-x64 && \
    chmod +x tailwindcss-linux-x64 && \
    mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss

# Copy dependency files
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build CSS
RUN tailwindcss -i ui/assets/input.css -o ui/static/css/output.css --minify

# Build the application with embedded files
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o main ./cmd/web

# Final stage
FROM debian:bookworm-slim

WORKDIR /app

# Install necessary runtime dependencies including SQLite
RUN apt-get update && apt-get install -y \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# Copy only the binary from builder
COPY --from=builder /app/main .

# Create directory for SQLite database
RUN mkdir -p ./data && \
    chmod 755 ./data

# Expose the production port
EXPOSE 4000

# Set environment variables
ENV PORT=4000

# Command to run migrations and start the application
CMD ["./main", "-addr", "0.0.0.0", "-dsn", "data/db.sqlite"]