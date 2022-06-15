# syntax=docker/dockerfile:1.2

FROM registry.suse.com/bci/nodejs:14

RUN zypper --non-interactive up && \
    zypper --non-interactive install nginx

RUN mkdir /app && \
    mkdir /etc/nginx/nginxconfig.io && \
    mkdir /srv/www/htdocs/skyscraper-web && \
    ln -sf /dev/stdout /var/log/nginx/access.log && \
    ln -sf /dev/stderr /var/log/nginx/error.log && \
    chown nginx:nginx /app && \
    chown nginx:nginx -R /srv/www/htdocs

COPY web/package*.json /app/
RUN cd /app && npm ci

COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY docker/config/* /etc/nginx/nginxconfig.io/
COPY docker/vhosts/* /etc/nginx/vhosts.d/
COPY docker/entrypoint.sh /usr/local/bin/docker-entrypoint.sh

COPY web /app

WORKDIR /srv/www/htdocs/skyscraper-web
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]

EXPOSE 8080
STOPSIGNAL SIGQUIT
USER nginx
CMD ["nginx", "-g", "daemon off;"]
