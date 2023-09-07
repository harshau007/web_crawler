FROM golang:alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o crawler

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/crawler .
ENTRYPOINT [ "./crawler" ]