FROM golang:alpine3.22 as builder

WORKDIR /app

COPY ./ ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goapp ./cmd/api

# == Run Stage ==

FROM alpine:3.22 as run

RUN apk update && apk add --no-cache \
    ca-certificates \
    tzdata 

ENV TZ=Asia/Bangkok

WORKDIR /app

COPY --from=builder /app/goapp ./goapp
COPY ./migrations ./migrations

EXPOSE 8080

CMD ["./goapp"]
