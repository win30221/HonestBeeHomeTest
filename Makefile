build:
	mkdir -p bin
	go build -o bin/client client.go
	go build -o bin/server server.go
	go build -o bin/externalapi externalapi.go