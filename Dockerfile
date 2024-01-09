FROM golang:latest
LABEL authors="a"

WORKDIR /cmd/

COPY . .

# Собираем приложение
RUN go build -o myapp cmd/main.go

# Определяем команду для запуска приложения при старте контейнера
CMD ["./cmd/myapp"]
