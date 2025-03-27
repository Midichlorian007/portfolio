FROM golang:1.23 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:2.21

COPY --from=builder /app/main .

CMD ["./main"]