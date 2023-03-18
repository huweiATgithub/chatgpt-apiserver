# syntax=docker/dockerfile:1.2
FROM golang:1.20-alpine3.17 AS builder
WORKDIR /apiserver
COPY . .
RUN go build -o /chatgpt-apiserver

FROM alpine:3.17
COPY --from=builder /chatgpt-apiserver /chatgpt-apiserver
WORKDIR /
CMD ["/chatgpt-apiserver"]
