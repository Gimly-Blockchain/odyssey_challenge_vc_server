# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Gianguido Sorà <me@gsora.xyz>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -installsuffix cgo -o dsbapi .

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

ARG LOG_DIR=/root/logs
ARG DSB_DIR=/root/dsb

# Create Log Directory
RUN mkdir -p ${LOG_DIR}
RUN mkdir -p ${DSB_DIR}

ENV DSB_STORAGE_PATH="/root/dsb"
ENV DSB_LOG_PATH="/root/logs/dsb.log"

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/dsbapi .

# Declare volumes to mount
VOLUME [${LOG_DIR}]
VOLUME [${DSB_DIR}]

# Command to run the executable
CMD ./main
