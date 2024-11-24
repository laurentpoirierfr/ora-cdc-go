# Build stage
FROM golang:1.22.5-alpine3.19 AS builder

ENV GOSUMDB="off"

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go test ./... -v --tags=unit
RUN go build -o main ./cmd/server/main.go

# Certs stage
FROM alpine:3.20 as certs


# Run stage final
FROM scratch

WORKDIR /app

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/main .
COPY ./cmd/server/*.yaml .

EXPOSE 8080

CMD [ "/app/main" ]

