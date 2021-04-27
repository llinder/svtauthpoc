build-client:
	go build -o bin/client cmd/client/main.go

build-server:
	go build -o bin/client cmd/server/main.go

build: build-client  build-server

run-server:
	CompileDaemon -build="go build -o bin/server cmd/server/main.go" -command ./bin/server

