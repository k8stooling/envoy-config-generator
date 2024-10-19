FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /envoy-config-generator .

FROM gcr.io/distroless/base
COPY --from=builder /envoy-config-generator /envoy-config-generator
USER 101:101
ENTRYPOINT ["/envoy-config-generator"]