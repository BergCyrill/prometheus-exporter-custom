FROM golang:1.19 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o exporter ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates git docker-cli
WORKDIR /root/
COPY --from=builder /app/exporter .
COPY config.yaml /etc/config/config.yaml
CMD ["./exporter"]
