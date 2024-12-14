FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o shorten-api .

FROM alpine:latest

RUN apk --no-cache add redis

COPY --from=builder /app/shorten-api /usr/local/bin/shorten-api

EXPOSE 8080

CMD ["shorten-api"]

