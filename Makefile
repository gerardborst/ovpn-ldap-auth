BINARY=ovpn-ldap-auth
GIT_COMMIT=$(shell git rev-parse HEAD)
HAS_NO_TAG=$(shell git describe --contains ${GIT_COMMIT} || true)
ifneq ($(strip ${HAS_NO_TAG}),)
VERSION_TAG=$(shell git describe --contains ${GIT_COMMIT})
endif

# If no git tag is set, fallback to 'DEVELOPMENT'
ifeq ($(strip ${VERSION_TAG}),)
VERSION_TAG := "DEVELOPMENT"
endif

COMMIT_HASH=$(git rev-parse --short=8 HEAD 2>/dev/null)
LDFLAGS=-ldflags "-s -w \
	-X main.CommitHash=${COMMIT_HASH} \
	-X main.BuildTime=`date --iso-8601=seconds` \
	-X main.VersionTag=${VERSION_TAG}"

all: build

clean:
	go clean
	rm -rf ./target || true

# Installs our project: copies binaries
install:
	go install ${LDFLAGS}

build:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	go build -o ${BINARY} ${LDFLAGS} cmd/${BINARY}/*.go

linux: clean
	GOOS=linux GOARCH=amd64 go build -o ./target/linux_amd64/${BINARY} ${LDFLAGS} cmd/${BINARY}/*.go
	cd ./target/linux_amd64/ && tar -czvf ${BINARY}.linux-amd64.tar.gz ${BINARY} && \
	sha256sum ${BINARY}.linux-amd64.tar.gz > sha256sum.txt

.PHONY: build
