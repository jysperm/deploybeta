FROM ubuntu:16.04

RUN echo "deb http://repo.mysql.com/apt/ubuntu/ xenial mysql-8.0" > /etc/apt/sources.list.d/mysql.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv A4A9406876FCBD3C456770C88C718D3B5072E1F5

RUN echo "deb http://ppa.launchpad.net/nginx/stable/ubuntu xenial main" > /etc/apt/sources.list.d/nginx.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv 8B3981E7A6852F782CC4951600A6F0A3C300EE8C

RUN { \
    echo mysql-community-server mysql-community-server/data-dir select ''; \
    echo mysql-community-server mysql-community-server/remove-test-db select false; \
    echo mysql-community-server mysql-community-server/root-pass password ''; \
    echo mysql-community-server mysql-community-server/re-root-pass password ''; \
  } | debconf-set-selections

RUN apt-get update && apt-get install -y supervisor mysql-community-client mysql-community-server-core nginx jq curl
RUN mkdir /var/run/mysqld

ADD supervisord.conf /etc/supervisord.conf
ADD control-agent nginx.conf.sh mysql.conf.sh /usr/bin/

VOLUME /var/lib/mysql

# Proxy port
EXPOSE 3306
# MySQL port
EXPOSE 3307

CMD ["supervisord", "-n", "-c", "/etc/supervisord.conf"]
