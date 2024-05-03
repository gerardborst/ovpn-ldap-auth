BINARY=ovpn-ldap-auth
$(eval VERSION_TAG = $(shell git describe 2>/dev/null | cut -f 1 -d '-' 2>/dev/null))

# If no git tag is set, fallback to 'DEVELOPMENT'
ifeq ($(strip ${VERSION_TAG}),)
VERSION_TAG := "DEVELOPMENT"
endif

COMMIT_HASH=`git rev-parse --short=8 HEAD 2>/dev/null`
# LDFLAGS=-ldflags "-s -w \
# 	-X src/git.gerardborst-ovpn.nl/ldap-auth/cmd/ldap-auth/main.CommitHash=${COMMIT_HASH} \
# 	-X src/git.gerardborst-ovpn.nl/ldap-auth/cmd/ldap-auth/main.VersionTag=${VERSION_TAG}"
LDFLAGS=-ldflags "-s -w \
	-X main.CommitHash=${COMMIT_HASH} \
	-X main.BuildTime=${BUILD_TIME} \
	-X main.VersionTag=${VERSION_TAG}"

all: build

clean:
	go clean
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm -rf ./target || true

release: clean linux

# Installs our project: copies binaries
install:
	go install ${LDFLAGS}

build:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	go build -o ${BINARY} ${LDFLAGS} cmd/ovpn-ldap-auth/*.go

linux:
	GOOS=linux GOARCH=amd64 go build -o ./target/linux_amd64/${BINARY} ${LDFLAGS} cmd/ovpn-ldap-auth/*.go

.PHONY: build
