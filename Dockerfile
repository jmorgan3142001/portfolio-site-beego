# Build the Go binary
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .
FROM alpine:latest
# Install CA certificates 
RUN apk add --no-cache ca-certificates
WORKDIR /root/
# Copy the binary from the builder
COPY --from=builder /app/main .
# Copy static files and views
COPY --from=builder /app/static ./static
COPY --from=builder /app/views ./views
COPY --from=builder /app/conf ./conf

EXPOSE 8080
CMD ["./main"]