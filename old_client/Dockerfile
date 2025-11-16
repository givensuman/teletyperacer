FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o client .

FROM alpine:3.22.1

RUN apk --no-cache add ca-certificates=20250619-r0
WORKDIR /root/

COPY --from=builder /app/client .

CMD ["./client"]
