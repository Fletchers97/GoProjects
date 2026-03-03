FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /collector-app ./cmd/collector

FROM alpine:latest

RUN apk add --no-cache ca-certificates libc6-compat
WORKDIR /

COPY --from=builder /collector-app /collector-app

COPY --from=builder /app/templates /templates

COPY --from=builder /app/config.json /config.json

RUN chmod +x /collector-app

ENTRYPOINT ["/collector-app"]