FROM golang:1.5.4-alpine

# golang alpine doesn't have ONBUILD, do it manually
COPY . /go/src/app
WORKDIR /go/src/app
RUN apk add --no-cache nodejs git && \
    go get -d -v && \
    go install -v

# install uchiwa-web bower package and cleanup
RUN npm install --production --unsafe-perm && \
    apk del --no-cache git

VOLUME /config

CMD ["/go/bin/app", "-c", "/config/config.json"]

EXPOSE 3000
