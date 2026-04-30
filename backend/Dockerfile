# -- Stage 1: Build --
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary for the API
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api-server ./cmd/api

# Build the Go binary for the Seeder
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o seeder ./cmd/seed

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o admin-seeder ./cmd/adminseed

# -- Stage 2: Run --
FROM alpine:latest  

WORKDIR /root/

# Copy the migrations folder
COPY migrations ./migrations

# Copy BOTH compiled binaries from the builder stage
COPY --from=builder /app/api-server .
COPY --from=builder /app/seeder .
COPY --from=builder /app/admin-seeder .

# Expose the port your Go app runs on
EXPOSE 8080

# Command to run the API executable by default
CMD ["./api-server"]