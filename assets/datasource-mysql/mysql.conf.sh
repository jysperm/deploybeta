#!/usr/bin/env bash

cat <<EOF
[mysqld]
server-id = ${SERVERID}
user = root
port = 3307
pid-file = /var/run/mysqld/mysqld.pid
socket = /var/run/mysqld/mysqld.sock
datadir = /var/lib/mysql
log-bin = mysql-bin
secure-file-priv= NULL

!includedir /etc/mysql/conf.d/
EOF
