# Build Stage
FROM golang:1.24.1-alpine3.20 AS builder
WORKDIR /app
COPY . . 
RUN go build -o main main.go
RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

# Run Stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main . 
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY app.env . 
COPY start.sh . 
COPY wait-for.sh . 
RUN chmod +x /app/wait-for.sh 
COPY db/migration ./migration

EXPOSE 8080

ENTRYPOINT ["/app/start.sh"]
