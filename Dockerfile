# Шаг 1: билдим go-бинарник
FROM golang:1.24.4 AS builder
WORKDIR /planner

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# Шаг 2: минимальный Ubuntu-образ
FROM ubuntu:minimal

WORKDIR /planner

COPY --from=builder /planner/app .
COPY --from=builder /planner/web ./web

ENV TODO_DBFILE=scheduler.db
ENV TODO_PORT=7540
ENV TODO_PASSWORD=123456

EXPOSE 7540

CMD ["./app"]