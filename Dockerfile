# Build
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build a static binary (easier for Cloud Run)
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Run
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/conf ./conf
COPY --from=builder /app/static ./static
COPY --from=builder /app/views ./views

# Expose
EXPOSE 8080
CMD ["./main"]