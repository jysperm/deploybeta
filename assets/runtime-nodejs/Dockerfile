FROM buildpack-deps:xenial-scm

RUN rm /bin/sh && ln -s /bin/bash /bin/sh

{{if .AptMirror}}
RUN sed -i 's/archive.ubuntu.com/{{.AptMirror}}/g' /etc/apt/sources.list
{{end}}

RUN apt-get update \
  && apt-get install -y curl \
  && apt-get -y autoclean

ENV NVM_DIR /usr/local/nvm

ENV NODE_VERSION {{.NodeVersion}}

RUN {{.ProxyCommand}} curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.2/install.sh | bash

RUN source $NVM_DIR/nvm.sh \
  && echo $NODE_VERSION > .nvmrc \
  && {{if .NvmMirror}} NVM_NODEJS_ORG_MIRROR={{.NvmMirror}}{{end}} nvm install $NODE_VERSION \
  && nvm alias default $NODE_VERSION \
  && nvm use default

ENV NODE_PATH $NVM_DIR/v$NODE_VERSION/lib/node_modules
ENV PATH $NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH


ENV YARN_VERSION 0.24.6

RUN set -ex \
  && for key in \
    6A010C5166006599AA17F08146C2130DFD2497F5 \
  ; do \
    gpg --keyserver pgp.mit.edu --recv-keys "$key" || \
    gpg --keyserver keyserver.pgp.com --recv-keys "$key" || \
    gpg --keyserver ha.pool.sks-keyservers.net --recv-keys "$key" ; \
  done \
  && {{.ProxyCommand}} curl -fSLO --compressed "https://yarnpkg.com/downloads/$YARN_VERSION/yarn-v$YARN_VERSION.tar.gz" \
  && {{.ProxyCommand}} curl -fSLO --compressed "https://yarnpkg.com/downloads/$YARN_VERSION/yarn-v$YARN_VERSION.tar.gz.asc" \
  && gpg --batch --verify yarn-v$YARN_VERSION.tar.gz.asc yarn-v$YARN_VERSION.tar.gz \
  && mkdir -p /opt/yarn \
  && tar -xzf yarn-v$YARN_VERSION.tar.gz -C /opt/yarn --strip-components=1 \
  && ln -s /opt/yarn/bin/yarn /usr/local/bin/yarn \
  && ln -s /opt/yarn/bin/yarn /usr/local/bin/yarnpkg \
  && rm yarn-v$YARN_VERSION.tar.gz.asc yarn-v$YARN_VERSION.tar.gz

RUN mkdir -p "/home/deploybeta/app"
WORKDIR "/home/deploybeta/app"
ADD . "/home/deploybeta/app"

{{if .NpmMirror}}
RUN npm config set registry {{.NpmMirror}}
{{end}}
{{if .UseYarn}}
RUN yarn install
{{- else}}
RUN npm install
{{- end}}

CMD ["npm", "start"]

EXPOSE 3000
