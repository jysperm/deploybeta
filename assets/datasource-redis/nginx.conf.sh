#!/usr/bin/env bash

cat <<EOF
load_module /usr/lib/nginx/modules/ngx_stream_module.so;

worker_processes 1;
daemon off;

events {
  worker_connections 65535;
}

stream {
  upstream redis_backend {
    server ${BACKEND_HOST};
  }

  server {
    listen 6379;
    proxy_pass redis_backend;
  }
}
EOF
