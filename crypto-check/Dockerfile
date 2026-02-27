FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . .

RUN go test -v ./...

RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM alpine:latest

RUN apk add --no-cache sqlite-libs ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config.json .
COPY --from=builder /app/templates ./templates

EXPOSE 8080

CMD ["./main"]