PACKAGE=github.com/RadioCheckerApp/api/request

VERSION=`git describe --tags --abbrev=0`
BUILD=`date -u +"%Y-%m-%dT%H:%M:%SZ"`
REVISION=`git rev-parse --short HEAD`

LDFLAGS=-ldflags "-s -w -X ${PACKAGE}.Version=${VERSION} -X ${PACKAGE}.Build=${BUILD} -X ${PACKAGE}.Revision=${REVISION}"

build:
	dep ensure
	env GOOS=linux go build ${LDFLAGS} -o ../bin/api-aws/meta meta/main.go
	env GOOS=linux go build ${LDFLAGS} -o ../bin/api-aws/stations stations/main.go
	env GOOS=linux go build ${LDFLAGS} -o ../bin/api-aws/tracks tracks/main.go
	env GOOS=linux go build ${LDFLAGS} -o ../bin/api-aws/search search/main.go
	env GOOS=linux go build ${LDFLAGS} -o ../bin/api-aws/tracks-create tracks-create/main.go