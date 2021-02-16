clean:
	rm -rf build

build_server: clean
	go mod vendor
	go build -ldflags "-s" -o build/go-remote src/main/main.go

install_server: build_server
	sudo cp ./build/go-remote /usr/local/bin/go-remote

uninstall_server:
	sudo rm /usr/local/bin/go-remote

build_command_executor: clean
	cd command-executor && go mod vendor
	go build -ldflags "-s" -o ../build/go-remote-command-executor command-executor/src/main/main.go

install_command_executor: build_command_executor
	sudo cp ./build/go-remote-command-executor /usr/local/bin/go-remote-command-executor

uninstall_command_executor:
	sudo rm /usr/local/bin/go-remote-command-executor