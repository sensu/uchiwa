#!/bin/bash

set -e

echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
  race=""
  # The race detector is broken on Alpine. That is #14481 (and #9918).
  # So disable it for now.
  if [ "${GOARCH}" = "amd64" ] && [ ! -f /etc/alpine-release ]; then
    race="-race"
  fi

  go test $race -coverprofile=profile.out -covermode=atomic "$d"
  if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
  fi
done

exit 0
