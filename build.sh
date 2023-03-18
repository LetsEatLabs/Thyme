#!/bin/bash
#

GOOS=darwin go build -o thyme-macos .
GOOS=linux go build -o thyme-linux .

s3cmd put thyme-macos s3://letseatlabs-data/bin/thyme/macos/thyme --acl-public
s3cmd put thyme-linux s3://letseatlabs-data/bin/thyme/linux/thyme --acl-public

rm thyme-macos
rm thyme-linux
