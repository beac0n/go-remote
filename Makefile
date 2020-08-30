FORCE: ;

build: FORCE
	go build -o build/client src/client/main.go

clean:
	rm -rf build
