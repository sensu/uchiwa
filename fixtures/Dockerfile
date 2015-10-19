FROM golang:1.5.1-onbuild

# update debian packages
RUN apt-get update

# install uchiwa-web bower package
RUN apt-get install -yq nodejs npm git wget
RUN ln -s /usr/bin/nodejs /usr/bin/node
RUN npm install --production --unsafe-perm

VOLUME /config

CMD ["/go/bin/app", "-c", "/config/config.json"]

EXPOSE 3000
