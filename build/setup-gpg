#!/bin/sh
set -e

rm -rf /home/travis/.gnupg

aws s3 cp s3://sensu-omnibus-cache/gpg/sensu-io-gpg.tar .
tar -xvf sensu-io-gpg.tar

cp .rpmmacros /home/travis/.rpmmacros
cp -R .gnupg /home/travis/.gnupg
