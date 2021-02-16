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
	cd command-executor && go build -ldflags "-s" -o ../build/go-remote-command-executor src/main/main.go

install_command_executor: build_command_executor
	sudo cp ./build/go-remote-command-executor /usr/local/bin/go-remote-command-executor

uninstall_command_executor:
	sudo rm /usr/local/bin/go-remote-command-executor

install: install_command_executor install_server

uninstall: uninstall_command_executor uninstall_server

test: build_command_executor
	-sudo killall go-remote-command-executor
	-rm /tmp/go-remote.sock
	mkdir -p "/tmp/go-remote"
	sudo ./build/go-remote-command-executor -user-name $$USER -command-timeout 1 -command-start 'install -m 777 /dev/null /tmp/go-remote/start' &
	go test src/main/main.go src/main/main_test.go src/main/test_utils.go
	sudo killall go-remote-command-executor
	rm /tmp/go-remote.sock
	rmdir "/tmp/go-remote"