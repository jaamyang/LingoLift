# Stage 1: Build the frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend

# Copy package files and install dependencies
COPY frontend/package*.json ./
RUN npm ci

# Copy frontend source and build
COPY frontend/ .
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app

# Copy go module files and download dependencies
COPY backend/go.mod backend/go.sum ./

RUN go mod download

# Copy backend source code
COPY backend/ .

# Copy the built frontend into the backend directory
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Stage 3: Create the final lightweight image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=backend-builder /app/server .

# Copy the frontend dist files
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Copy existing models and wasm files
COPY models/ ./frontend/dist/models/
COPY wasm/ ./frontend/dist/wasm/

# Expose port 8080
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the server
CMD ["./server"]
