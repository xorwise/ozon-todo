FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN go build -o main .

FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite

WORKDIR /app
 
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
