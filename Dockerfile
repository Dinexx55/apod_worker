FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

RUN ls -la main

FROM debian:bookworm AS runner

WORKDIR /usr/app

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]