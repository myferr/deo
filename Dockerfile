FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY . .

RUN go build -o /deo

FROM alpine:latest

WORKDIR /app

COPY --from=builder /deo /deo
COPY --from=builder /app/templates /app/templates

EXPOSE 6741

CMD ["sh", "-c", "cd /app && /deo"]
