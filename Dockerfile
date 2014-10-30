FROM golang:1.3.3

# install debian packages
RUN apt-get update && apt-get install -yq nodejs npm git wget

RUN /usr/src/go/bin/go get github.com/tools/godep

RUN ln -s /usr/bin/nodejs /usr/bin/node

ADD . /src

WORKDIR /src

RUN npm install --production --unsafe-perm

RUN mv /src/docker/start /start && chmod 0755 /start

RUN go get -d -v

RUN go build -o uchiwa_binary uchiwa.go

VOLUME /config

EXPOSE 3000
CMD ["/start"]
