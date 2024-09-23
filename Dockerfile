FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/proxy

RUN go build -o /app/proxyapp main.go

FROM alpine:latest

ARG APP_ENV

RUN echo "APP_ENV=${APP_ENV}" > .env

COPY --from=builder /app/proxyapp /proxyapp
COPY --from=builder /app/config /config

EXPOSE 8080
EXPOSE 3333

# Command to run the executable
CMD ["./proxyapp"]