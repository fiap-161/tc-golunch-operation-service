# Build stage

FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server ./cmd/api/main.go

# Final stage
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/server /server
COPY --from=builder /app/conf /conf

ENTRYPOINT ["/server"]
