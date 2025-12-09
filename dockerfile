# Build stage
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o golab ./cmd/golab

# Final minimal image
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/golab .
CMD ["./golab"]
