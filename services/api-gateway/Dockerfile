FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the entire project to have access to go.work and pkg
COPY . .

# Build the specific service
RUN go build -o /server services/api-gateway/cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /server .

EXPOSE 8080

CMD ["./server"]
