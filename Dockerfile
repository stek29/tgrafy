# Builder image
FROM golang:1.12.0-alpine3.9 AS builder

RUN apk add git

WORKDIR /build
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

ADD . .
RUN go build -o tgrafy .

# Final image
FROM alpine:3.9

# Copy all binaries
COPY --from=builder /build/tgrafy /tgrafy

EXPOSE 80
ENTRYPOINT [ "/tgrafy" ]
