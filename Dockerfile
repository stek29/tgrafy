# Builder image
FROM golang:1.12.0-alpine3.9 AS builder

# go mod needs git
RUN apk add git

WORKDIR /build

# Go Mod init
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
ADD . .

# Build the executable
RUN go build -o tgrafy .


# Final image
FROM alpine:3.9

# Install CA Certs for TLS
RUN apk add --no-cache ca-certificates

# Copy all binaries
COPY --from=builder /build/tgrafy /tgrafy

# Default listen port is 80
EXPOSE 80
ENTRYPOINT [ "/tgrafy" ]
