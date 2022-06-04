#!/usr/bin/env bash

cd /app && \
  ./node_modules/.bin/ng build --configuration production && \
  cp -rf /app/dist/skyscraper-web/* /srv/www/htdocs/skyscraper-web && \
  chown nginx:nginx -R /srv/www/htdocs/skyscraper-web

exec "$@"
