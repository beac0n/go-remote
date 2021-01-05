FORCE: ;

build: FORCE
	go build -ldflags "-s" -o build/go-remote src/main/main.go

clean:
	rm -rf build
