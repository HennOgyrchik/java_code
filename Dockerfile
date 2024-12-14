FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY . .
RUN  mkdir /app
RUN go build -o /app/java_code ./cmd/

FROM alpine:latest
COPY --from=builder /app /app
CMD ["./app/java_code"]