FROM golang:alpine

# Download and install the latest Goose CLI
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Set Goose as the default command when this container boots
ENTRYPOINT ["goose"]