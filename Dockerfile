FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /server ./cmd/server

FROM alpine:3.18
COPY --from=builder /server /server
COPY static /static
EXPOSE 8080
ENTRYPOINT ["/server"]
