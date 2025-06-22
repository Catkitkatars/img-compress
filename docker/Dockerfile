FROM alpine:3.22

RUN apk update && apk add --no-cache \
    go \
    nginx \
    sqlite \
    supervisor \
    bash

WORKDIR /app

COPY . .

WORKDIR /app/cmd/img-compress
RUN go mod init img-compress
RUN go mod tidy

RUN go build -o /app/img-compress

WORKDIR /app

COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY docker/conf.d/ /etc/nginx/conf.d/
COPY docker/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

RUN chmod +x /app/img-compress

EXPOSE 8080

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]