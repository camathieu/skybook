# Stage 1: Build the frontend (Vue 3 / Vite)
FROM node:22-alpine AS ui-builder
WORKDIR /webapp
COPY webapp/package*.json ./
RUN npm ci
COPY webapp/ ./
RUN npm run build

# Stage 2: Build the Go server (combines with frontend dist)
FROM golang:1.25-alpine AS go-builder
# Install gcc and musl-dev to compile go-sqlite3 correctly with CGO
RUN apk add --no-cache gcc musl-dev
WORKDIR /server

# Cache Go modules
COPY server/go.mod server/go.sum ./
RUN go mod download

# Prepare the build workspace
COPY server/ ./
# Create the directory for the frontend files expected by go:embed
RUN mkdir -p server/dist
# Copy the built UI dist from Stage 1 into the location main.go embed expects
COPY --from=ui-builder /webapp/dist ./server/dist/

# Build statically linked binary with CGO enabled (required for sqlite)
RUN CGO_ENABLED=1 \
    go build -tags "osusergo,netgo,sqlite_omit_load_extension" \
    -ldflags "-linkmode external -extldflags -static" \
    -o skybook .

# Stage 3: Minimal runtime image
FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

# Copy the static binary from builder
COPY --from=go-builder /server/skybook .

# Application uses port 8080 by default
EXPOSE 8080

# The sqlite database should live on a persistent volume
ENV SKYBOOK_DATABASE_PATH=/data/skybook.db
VOLUME /data

# Run the server
ENTRYPOINT ["/app/skybook", "serve"]
