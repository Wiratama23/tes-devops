# ==========================================
# STAGE 1: The Heavy Builder
# ==========================================
FROM golang:alpine AS builder

# Download and install the latest Goose CLI
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# ==========================================
# STAGE 2: The Tiny Final Image
# ==========================================
FROM alpine:latest

# Install basic certificates so Goose can securely connect to PostgreSQL
RUN apk --no-cache add ca-certificates

WORKDIR /migrations

# Copy ONLY the finished, compiled binary from Stage 1
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Bake your SQL files directly into the image
COPY ./migrations/*.sql /migrations/

# Set Goose as the default command
ENTRYPOINT ["goose"]