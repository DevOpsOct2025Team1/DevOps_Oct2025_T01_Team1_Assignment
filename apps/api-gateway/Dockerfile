FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY apps/ ./apps
COPY go.work ./
COPY go.work.sum ./

WORKDIR /build/apps/api-gateway
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/api-gateway ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates netcat-openbsd

WORKDIR /app

COPY --from=builder /app/api-gateway .

EXPOSE 3000

CMD ["./api-gateway"]