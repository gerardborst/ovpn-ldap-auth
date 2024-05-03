#!/bin/sh

git describe --exact-match --tags
if [ "$?" -ne 0 ]; then
    echo "No tag on commit, version will be beta"
    VERSION=beta
else
    VERSION=$(git describe --exact-match --tags $(git rev-list --tags --max-count=1))
fi
COMMIT_HASH=`git rev-parse --short=8 HEAD 2>/dev/null`

echo "Version is: ${VERSION:-beta}"
echo "Commit hash is: ${COMMIT_HASH}"

# APP_NAME="ldap-auth"

rm -fr target/*

gox -ldflags "-X main.VersionTag=${VERSION:-beta} \
    -X main.CommitHash=$COMMIT_HASH" \
    -os="linux darwin" -arch="amd64" -output="target/{{.OS}}-{{.Arch}}/ldap-auth" git.gerardborst-ovpn.nl/ldap-auth/cmd/ldap-auth
    
tar -czf target/ldap-auth-${VERSION}.linux-amd64.tar.gz -C target/linux-amd64 ldap-auth
tar -czf target/ldap-auth-${VERSION}.darwin-amd64.tar.gz -C target/darwin-amd64 ldap-auth