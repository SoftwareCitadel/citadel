# Build stage
FROM golang:alpine AS builder

WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build -p 4 --ldflags "-extldflags -static" -o web ./cmd/web

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/web /app/

ENTRYPOINT ["/app/web", "serve"]
