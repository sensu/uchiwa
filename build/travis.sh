#!/bin/bash

set -e

COMMIT=$(git rev-parse HEAD)

if [ "$(git describe --tags --exact-match "$COMMIT")" ]; then
  TAGS=$(git tag -l --points-at HEAD)
  TAG=$(echo "$TAGS" | sort -r | head -1)

  echo "======================== Found tags $TAGS on commit $COMMIT"
  echo "======================== Selected $TAG as latest tag"
  PACKAGE_VERSION=$(echo "$TAG" | awk -F'-' '{print $1}')
  BUILD_NUMBER=$(echo "$TAG" | awk -F'-' '{print $2}')

  if [ -z "$PACKAGE_VERSION" ] || [ -z "$BUILD_NUMBER"  ]; then
    echo "The tag '$TAG' does not contain any valid version (x.y.z-a)"
    exit 2
  fi

  echo "======================== Running tests"
  GOARCH=$GOARCH ./build/tests.sh

  echo "======================== Prepare RPM signing"
   pip install --user awscli
  export PATH=$PATH:$HOME/.local/bin
  build/setup-gpg

  echo "======================== Building the packages"
  PACKAGE_VERSION=$PACKAGE_VERSION BUILD_NUMBER=$BUILD_NUMBER \
  GOOS=$GOOS GOARCH=$GOARCH \
  rake -f build/Rakefile

  exit
else
  echo "Commit ${COMMIT} is not tagged, running tests"
  GOARCH=$GOARCH ./build/tests.sh
fi
