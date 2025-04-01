FROM nginx:latest

# Копируем конфиг
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

# Копируем файлы фронтенда
COPY frontend/ /usr/share/nginx/html

EXPOSE 80
