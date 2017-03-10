#!/bin/bash

set -e

echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
  race=""
  if [ $GOARCH == "amd64" ]; then
    race="-race"
  fi

  go test $race -coverprofile=profile.out -covermode=atomic $d
  if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
  fi
done

exit 0
