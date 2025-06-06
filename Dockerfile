FROM golang:1.24.3-alpine3.21 as builder

WORKDIR /app

COPY ./ ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goapp ./main.go

# == Run Stage ==

FROM alpine:3.21 as run

RUN apk update && apk add --no-cache \
        ca-certificates \
        tzdata 

ENV TZ=Asia/Bangkok

WORKDIR /app

COPY --from=builder /app/goapp ./goapp

EXPOSE 8080

CMD ["./goapp"]
