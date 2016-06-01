FROM golang:1.5.4-alpine

# golang alpine doesn't have ONBUILD, do it manually, then run npm and cleanup
COPY . /go/src/app
WORKDIR /go/src/app
RUN apk add --no-cache nodejs git && \
    go get -d -v && \
    go install -v && \
    npm install --production --unsafe-perm && \
    npm dedupe && \
    apk del --no-cache git

VOLUME /config

CMD ["/go/bin/app", "-c", "/config/config.json"]

EXPOSE 3000
