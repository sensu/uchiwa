FROM golang:1.9.2-alpine

# golang alpine doesn't have ONBUILD, do it manually, then run npm and cleanup
COPY . /go/src/github.com/sensu/uchiwa
WORKDIR /go/src/github.com/sensu/uchiwa
RUN apk add --no-cache nodejs-npm git && \
    go install -v && \
    npm install --production --unsafe-perm && \
    npm dedupe && \
    apk del --no-cache git nodejs-npm && \
    rm -rf /go/src/github.com/sensu/uchiwa/node_modules

VOLUME /config

CMD ["/go/bin/uchiwa", "-c", "/config/config.json"]

EXPOSE 3000
