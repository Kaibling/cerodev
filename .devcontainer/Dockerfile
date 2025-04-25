# Use an official Go runtime as a parent image
FROM golang:1.24.2-bookworm

# Set the Current Working Directory inside the container
WORKDIR /workspace
# Set environment variables for Go
ENV GOPATH=/go
ENV GOROOT=/usr/local/go
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

CMD tail -f /dev/null