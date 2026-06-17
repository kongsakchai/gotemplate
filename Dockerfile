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

# Copy binary app and migrations file
COPY --from=builder /app/goapp ./goapp
COPY ./migrations ./migrations

# Permissions
RUN chown -R app:app /app && chmod 0755 /app/goapp

# Drop root
USER app:app

EXPOSE 8080

CMD ["./goapp"]
