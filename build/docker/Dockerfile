FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/notification-service /app/cmd/main.go


FROM alpine:3 AS app

WORKDIR /app

COPY --from=builder /app/notification-service /app/notification-service
COPY --from=builder /app/config /app/config

CMD ["/app/notification-service"]