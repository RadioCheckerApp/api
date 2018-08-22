build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o ../bin/api-aws/meta meta/main.go
	env GOOS=linux go build -ldflags="-s -w" -o ../bin/api-aws/stations stations/main.go
	env GOOS=linux go build -ldflags="-s -w" -o ../bin/api-aws/tracks tracks/main.go
	env GOOS=linux go build -ldflags="-s -w" -o ../bin/api-aws/search search/main.go
	env GOOS=linux go build -ldflags="-s -w" -o ../bin/api-aws/tracks-create tracks-create/main.go