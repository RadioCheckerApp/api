build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o ../bin/api-aws/meta meta/main.go

test:
	dep ensure
	go test -race ./*/