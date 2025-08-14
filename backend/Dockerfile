# Stage 1: Build the frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
# Install git to clone the model repo
RUN apk --no-cache add git git-lfs
COPY frontend/package*.json ./
RUN npm install

# Copy wasm files from node_modules into a public directory so vite includes them in the build output.
RUN mkdir -p public/wasm && cp -r node_modules/@huggingface/transformers/dist/wasm/* public/wasm/

# Clone the model repository from Hugging Face Hub
RUN git clone https://huggingface.co/Xenova/whisper-tiny public/models/Xenova/whisper-tiny

COPY frontend/ .
RUN npm run build

# Stage 2: Build the Go backend with embedded frontend
FROM golang:1.21 AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Copy go source
COPY backend-go/ .
# Copy the built frontend into the backend-go directory so it can be embedded.
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o /server .

# Stage 3: Create the final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the Go binary from the go-builder stage
COPY --from=go-builder /server .

EXPOSE 8080
CMD ["./server"] 