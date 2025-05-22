# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24.3-alpine3.21 AS builder
ARG CGO_ENABLED=0
ARG GOOS=linux
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /app/ruuvi-gateway-prometheus

# Test stage
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Runtime image
FROM scratch
COPY --from=builder /app/ruuvi-gateway-prometheus /ruuvi-gateway-prometheus
ENTRYPOINT ["/ruuvi-gateway-prometheus"]
EXPOSE 9090
LABEL org.opencontainers.image.source=https://github.com/jkjuopperi/ruuvi-gateway-prometheus
LABEL org.opencontainers.image.description="Ruuvi Gateway Prometheus Exporter"
LABEL org.opencontainers.image.licenses=BSD-2-Clause
