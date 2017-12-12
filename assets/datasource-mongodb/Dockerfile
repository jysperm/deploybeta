FROM ubuntu:16.04

RUN echo "deb [ arch=amd64,arm64 ] http://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/3.6 multiverse" > /etc/apt/sources.list.d/mongodb.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv 2930ADAE8CAF5059EE73BB4B58712A2291FA4AD5

RUN echo "deb http://ppa.launchpad.net/nginx/stable/ubuntu xenial main" > /etc/apt/sources.list.d/nginx.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv 8B3981E7A6852F782CC4951600A6F0A3C300EE8C

RUN apt-get update && apt-get install -y supervisor mongodb-org nginx jq curl

ADD supervisord.conf /etc/supervisord.conf
ADD control-agent nginx.conf.sh /usr/bin/

VOLUME /var/lib/mongodb

# Proxy port
EXPOSE 27017
# MongoDB port
EXPOSE 27018

CMD ["supervisord", "-n", "-c", "/etc/supervisord.conf"]
