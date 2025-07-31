# syntax=docker/dockerfile:1

FROM golang:tip-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o csv-processor ./cmd/processor

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/csv-processor .
COPY --from=builder /app/internal ./internal
COPY --from=builder /app/cmd ./cmd
COPY --from=builder /app/go.mod .
COPY --from=builder /app/go.sum .
ENTRYPOINT ["/app/csv-processor"]