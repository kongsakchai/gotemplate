FROM golang:1.24.3-alpine3.21 as builder

WORKDIR /app

COPY ./ ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./main.go

FROM alpine:3.21 as run
WORKDIR /app
COPY --from=builder /app/app .

CMD ["./app"]
