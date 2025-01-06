FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o fiber-api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/fiber-api .
ENV env=prod
ENV region=us-east-1
EXPOSE 8080
CMD ["./fiber-api"]
