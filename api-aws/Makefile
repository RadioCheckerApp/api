build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o ../bin/api-aws/meta meta/main.go
	env GOOS=linux go build -ldflags="-s -w" -o ../bin/api-aws/stations stations/main.go