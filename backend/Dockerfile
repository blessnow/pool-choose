# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server

# Production stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/data ./data
EXPOSE 8080
ENV DB_PATH=/app/data/cycle_stock.db
ENV PASSWORD=dayuchi
CMD ["./server"]