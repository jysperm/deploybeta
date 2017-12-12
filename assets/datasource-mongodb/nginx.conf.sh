#!/usr/bin/env bash

cat <<EOF
load_module /usr/lib/nginx/modules/ngx_stream_module.so;

worker_processes 1;
daemon off;

events {
  worker_connections 65535;
}

stream {
  upstream mongodb_backend {
    server ${BACKEND_HOST};
  }

  server {
    listen 27017;
    proxy_pass mongodb_backend;
  }
}
EOF
