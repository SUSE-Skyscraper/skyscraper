# syntax=docker/dockerfile:1.2

FROM registry.suse.com/bci/nodejs:14 as builder

RUN mkdir /app
WORKDIR /app

COPY package*.json /app/
RUN npm ci
COPY . .

RUN ./node_modules/.bin/ng build

FROM registry.suse.com/bci/bci-base:latest

RUN zypper --non-interactive dup && \
    zypper --non-interactive install nginx

RUN mkdir /etc/nginx/nginxconfig.io

COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY docker/config/* /etc/nginx/nginxconfig.io/
COPY docker/vhosts/* /etc/nginx/vhosts.d/

RUN ln -sf /dev/stdout /var/log/nginx/access.log \
    && ln -sf /dev/stderr /var/log/nginx/error.log \
    && mkdir /srv/www/htdocs/skyscraper-web

COPY --from=builder /app/dist/skyscraper-web /srv/www/htdocs/skyscraper-web

WORKDIR /srv/www/htdocs/skyscraper-web

EXPOSE 8080

STOPSIGNAL SIGQUIT

ENTRYPOINT ["nginx", "-g", "daemon off;"]
