# Используйте официальный образ Golang как базовый
FROM golang:1.18 as builder

# Установите рабочий каталог внутри контейнера
WORKDIR /app

# Копируйте go.mod и go.sum для кэширования зависимостей
COPY go.mod .
COPY go.sum .

# Скачивайте зависимости. Возможно, это улучшит повторное использование слоя и ускорит сборку
RUN go mod download

# Копируйте исходный код проекта в контейнер
COPY . .

# Собирайте ваше приложение. Указывайте путь до вашего main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/portfolio/

# Используйте небольшой образ для запуска приложения
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируйте собранное приложение из предыдущего шага
COPY --from=builder /app/main .

# Копируйте статические файлы, если они нужны приложению
COPY --from=builder /app/icons /root/icons
COPY --from=builder /app/media /root/media
COPY --from=builder /app/sectionsData.json .

# Укажите порт, на котором будет работать приложение
EXPOSE 8080

# Запускайте приложение
CMD ["./main"]


