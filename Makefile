clean:
	rm -rf build

build: clean
	go build -ldflags "-s" -o build/go-remote src/main/main.go

install: build
	sudo cp ./build/go-remote /usr/local/bin/go-remote

uninstall:
	sudo rm /usr/local/bin/go-remote