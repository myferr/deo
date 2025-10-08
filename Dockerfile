FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY . .

RUN go build -o /deo

FROM alpine:latest

COPY --from=builder /deo /deo

EXPOSE 6741

CMD ["/deo"]
