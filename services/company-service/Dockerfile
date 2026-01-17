FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o /server services/company-service/cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /server .

EXPOSE 8082

CMD ["./server"]
