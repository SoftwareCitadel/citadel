# Build stage
FROM golang:alpine AS builder

WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build -p 4 --ldflags "-extldflags -static" -o smtp ./cmd/smtp

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/smtp /app/

ENTRYPOINT ["/app/smtp"]
