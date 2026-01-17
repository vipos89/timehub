FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o /server services/report-service/cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /server .

EXPOSE 8085

CMD ["./server"]
