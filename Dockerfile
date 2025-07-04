# ---------- Stage 1: Build and test support ----------
    FROM golang:1.23.4-alpine AS builder

    # Install useful tools
    RUN apk add --no-cache curl iputils postgresql-client git
    
    # Set working directory
    WORKDIR /app
    
    # Copy go.mod and go.sum for dependency resolution
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy the entire source code (including tests)
    COPY . .
    
    # Optional: Run unit tests (uncomment to run during build)
    # RUN go test ./tests/...

    # Build the application
    RUN go build -o match3-app .
    
    # ---------- Stage 2: Lightweight runtime ----------
    FROM alpine:latest
    
    # Install CA certificates (required for HTTPS, etc.)
    RUN apk add --no-cache ca-certificates
    
    # Set working directory
    WORKDIR /app
    
    # Copy the built binary from the builder stage
    COPY --from=builder /app/match3-app .

    # Copy Swagger JSON and ReDoc HTML
    COPY ./redoc.html ./redoc.html
    COPY ./docs/swagger.json ./docs/swagger.json
    
    # Expose the app port
    EXPOSE 8080
    
    # Default command to run the app
    CMD ["./match3-app"]
    