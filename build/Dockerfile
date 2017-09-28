FROM golang:1.8.1-alpine

ENV GOARCH=amd64

# Disable cgo so that binaries we build will be fully static.
ENV CGO_ENABLED=0

RUN apk add --no-cache \
  bash \
  build-base \
  git \
  libffi-dev \
  make \
  nodejs \
  rpm \
  ruby-dev \
  ruby-rake \
  tar && \
  gem install rake -v "10.5.0" --no-ri --no-rdoc && \
  gem install fpm -v "1.8.1" --no-ri --no-rdoc && \
  rm -rf /usr/lib/ruby/gems/*/cache/* && \
  rm -rf ~/.gem
