# Use an offical golang image to create the binary.
# https://hub.docker.com/_/golang
FROM 1.20.2-alpine3.17 AS builder

# Create a working directory in the image.
WORKDIR /app

# Copy local go.mod and, if present, go.sum to the container image.
COPY go.* ./
# Download necessary Go modules; https://go.dev/ref/mod#go-mod-download
RUN go mod download

# Copy local code to the container image.
COPY . ./

#Build the binary.
RUN go build -v -o /server




  # Use the official Debian slim image for a lean production container.
  # https://hub.docker.com/_/debian
  # https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-buildsFROM gcr.io/distroless/base-debian10
FROM debian:buster-slim

WORKDIR /

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server

EXPOSE 8080

USER nonroot:nonroot

# Run the web service on container startup.
CMD [ "/app/server" ]
