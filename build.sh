#!/bin/sh

VERSION="DEVELOPMENT"
COMMIT_HASH=`git rev-parse --short=8 HEAD 2>/dev/null`

echo "VERSION is ${VERSION}"
echo "COMMIT_HASH is ${COMMIT_HASH}\n"

APP_NAME="ldap-auth"

if [ -f $APP_NAME ] ; then
    rm $APP_NAME
fi

go build -o $APP_NAME \
    -ldflags "-X main.VersionTag=$VERSION -X main.CommitHash=$COMMIT_HASH" \
    cmd/ldap-auth/*.go