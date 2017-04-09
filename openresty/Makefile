install:
	apt-get update
	apt-get -y install libreadline-dev libncurses5-dev libpcre3-dev libssl-dev perl make build-essential curl

	mkdir openresty

	wget https://openresty.org/download/openresty-1.11.2.2.tar.gz -O openresty.tar.gz
	tar -zxf openresty.tar.gz -C openresty --strip-components=1 && rm openresty.tar.gz

	git clone https://github.com/yzprofile/ngx_http_dyups_module.git ngx_http_dyups_module
	git clone https://github.com/pintsized/lua-resty-http.git lua-resty-http
	cp -r lua-resty-http/lib/resty lib

	cd openresty && ./configure --add-module=../ngx_http_dyups_module --with-http_stub_status_module
	DESTDIR=`pwd` make -C openresty install

	rm -r openresty ngx_http_dyups_module lua-resty-http

run:
	./usr/local/openresty/nginx/sbin/nginx -p `pwd` -c conf/nginx.conf
