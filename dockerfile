FROM nginx:1.31.5

COPY nginx.conf /etc/nginx/nginx.conf
COPY dist /usr/share/nginx/html/
