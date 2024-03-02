# Используйте базовый образ Nginx
FROM nginx:latest

# Копирование кастомных конфигурационных файлов Nginx
COPY nginx.conf /etc/nginx/nginx.conf 
COPY ./default.conf /etc/nginx/conf.d/default.conf

# Использование томов для динамического контента
VOLUME ["/usr/share/nginx/html/media", "/usr/share/nginx/html/icons"]

# Экспозиция порта
EXPOSE 80

# Запуск Nginx
CMD ["nginx", "-g", "daemon off;"]
