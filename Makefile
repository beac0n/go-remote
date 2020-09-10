FORCE: ;

build: build_client build_server

build_client: FORCE
	go build -o build/go-remote-client src/client/main.go

build_server: FORCE
	go build -o build/go-remote-server src/server/main.go

clean:
	rm -rf build
