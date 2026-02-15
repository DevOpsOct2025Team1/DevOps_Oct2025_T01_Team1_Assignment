FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY apps/ ./apps
COPY go.work ./
COPY go.work.sum ./

WORKDIR /build/apps/user-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/user-service ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates netcat-openbsd && \
    update-ca-certificates

WORKDIR /app

COPY --from=builder /app/user-service .

EXPOSE 8080

CMD ["./user-service"]