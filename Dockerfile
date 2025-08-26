# Шаг 1: билдим go-бинарник
FROM golang:1.24 AS builder
WORKDIR /planner

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# Шаг 2: минимальный Ubuntu-образ
FROM ubuntu:22.04

WORKDIR /planner

COPY --from=builder /planner/app .
COPY --from=builder /planner/web ./web

EXPOSE 7540

CMD ["./app"]