#!/usr/bin/env bash

cat <<EOF
load_module /usr/lib/nginx/modules/ngx_stream_module.so;

worker_processes 1;
daemon off;

events {
  worker_connections 65535;
}

stream {
  upstream mysql_backend {
    server ${BACKEND_HOST};
  }

  server {
    listen 3306;
    proxy_pass mysql_backend;
  }
}
EOF
