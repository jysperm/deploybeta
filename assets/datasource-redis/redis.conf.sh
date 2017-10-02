#!/usr/bin/env bash

cat <<EOF
port 6380

save 900 1
save 300 10
save 60 10000

protected-mode no
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
EOF

if [ $ROLE = "slave" ]; then
  cat <<EOF
slaveof ${MASTER_HOST//:/ }
EOF
fi
