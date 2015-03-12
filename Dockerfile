FROM golang:1.3.3-onbuild

# install debian packages
RUN apt-get update \
&& apt-get install -yq nodejs npm git wget \
&& ln -s /usr/bin/nodejs /usr/bin/node \
&& npm install --production --unsafe-perm \
&& mv ./docker/start /start && chmod 0755 /start

VOLUME /config

CMD ["/start"]

EXPOSE 3000
