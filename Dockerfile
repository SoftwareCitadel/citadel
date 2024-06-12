FROM golang:alpine AS builder

WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build --ldflags "-extldflags -static" -o caesar .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/caesar /app/

ENTRYPOINT ["/app/caesar", "--serve"]
