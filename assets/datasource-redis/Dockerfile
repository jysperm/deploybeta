FROM ubuntu:16.04

RUN echo "deb http://ppa.launchpad.net/chris-lea/redis-server/ubuntu xenial main" > /etc/apt/sources.list.d/redis-server.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv 136221EE520DDFAF0A905689B9316A7BC7917B12

RUN echo "deb http://ppa.launchpad.net/nginx/stable/ubuntu xenial main" > /etc/apt/sources.list.d/nginx.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv 8B3981E7A6852F782CC4951600A6F0A3C300EE8C

RUN apt-get update && apt-get install -y supervisor redis-server nginx jq curl

ADD supervisord.conf /etc/supervisord.conf
ADD control-agent nginx.conf.sh redis.conf.sh /usr/bin/

VOLUME /var/lib/redis

# Proxy port
EXPOSE 6379
# Redis port
EXPOSE 6380

CMD ["supervisord", "-n", "-c", "/etc/supervisord.conf"]
