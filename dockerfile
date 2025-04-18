# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o flight-service ./cmd/flight-service

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/flight-service .
COPY .env .env
EXPOSE 3000
CMD ["./flight-service"]
