FROM golang:1.23

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем код приложения
COPY backend/ .

# Собираем Go-приложение
RUN go build -o server

# Открываем порт 5000
EXPOSE 5000

# Запускаем сервер
CMD ["/app/server"]
