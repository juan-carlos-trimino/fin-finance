# To test the image using docker:
# Build the image.
# $ docker build -t webserver:1.0.0 .
# List the images.
# $ docker image ls
# To remove an image.
# docker rmi <image-id>
# To start a container.
# docker run -d --name finances -p 8000:8000 webserver
# To list all running containers.
# docker ps
# To stop a running container.
# $ docker stop <container-id>
# To run commands inside an image.
# Alpine images provide the Almquist shell (ash) from BusyBox.
# $ docker exec -it <container-id> ash
# -------------------------------------------------------------------------------------------------
# Use an offical golang image to create the binary.
# Alpine images provide the Almquist shell (ash) from BusyBox.
FROM golang:1.20.2-alpine3.17 AS builder

# Create a working directory in the image.
WORKDIR /goapp/

# Copy go.mod and, if present, go.sum from the local machine to the container image.
COPY ./src/go.* ./
# Copy the code from the local machine to the container image.
COPY ./src/main.go ./
COPY ./src/finances ./finances/
COPY ./src/mathutil ./mathutil/
COPY ./src/webfinances ./webfinances/

# Download necessary Go modules; https://go.dev/ref/mod#go-mod-download
RUN go mod download

# Build the binary; https://go.googlesource.com/sublime-build/+/HEAD/docs/configuration.md
RUN go build -o godir/webserver -v main.go

FROM alpine:3.17

WORKDIR /

# Copy the binary to the production image from the builder stage.
COPY --from=builder godir/ goapp/

EXPOSE 8000

USER 1000:1000

# Run the web service on container startup.
ENTRYPOINT ["./goapp/webserver"]
