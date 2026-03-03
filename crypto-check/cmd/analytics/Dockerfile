FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /analytics-app ./cmd/analytics


FROM alpine:latest

RUN apk add --no-cache ca-certificates libc6-compat
WORKDIR /root/


COPY --from=builder /analytics-app /analytics-app


RUN chmod +x /analytics-app


ENTRYPOINT ["/analytics-app"]