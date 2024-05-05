# Используйте официальный образ Golang как базовый
FROM golang:1.22.1

WORKDIR /app

# Копируйте все файлы из текущего каталога в контейнер
COPY . .

# Скачивание зависимостей
RUN go mod tidy

# Соберите приложение
RUN go build -o server ./cmd/main

# Запустите приложение
CMD ["./server"]