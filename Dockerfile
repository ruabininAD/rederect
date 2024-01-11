#  sudo docker build -t redirect-img . &&  sudo docker run -p 3306:3306 redirect-img
# работает  sudo docker build -t redirect-img . &&  sudo docker run  redirect-img
# Используем образ для запуска приложений Go
FROM golang:alpine AS builder

# Копируем исходный код в образ
WORKDIR /cmd/
COPY . .



# Компилируем приложение
RUN go build -o myapp cmd/main.go

# Создаем Docker-образ с минимальным размером
FROM alpine:latest

# Копируем исполняемый файл из предыдущего образа
COPY --from=builder /cmd/myapp /usr/bin/myapp
COPY --from=builder /cmd/config.yaml /usr/bin/config.yaml

EXPOSE 8082
EXPOSE 2112
EXPOSE 3306

# Запускаем приложение
CMD ["myapp"]