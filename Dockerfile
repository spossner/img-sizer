# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

ARG GOOS=linux
ARG GOARCH=amd64

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY cmd/img-sizer/ ./cmd/img-sizer/
COPY internal/ ./internal/
COPY config/ ./config/

# Build the application
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o img-sizer ./cmd/img-sizer

# Final stage
FROM gcr.io/distroless/static-debian11

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/img-sizer .
COPY --from=builder /app/config ./config

# Set default environment
ENV APP_ENV=prod
ENV PORT=8080

# Expose port
EXPOSE 8080

# Run the application
CMD ["./img-sizer"] 
