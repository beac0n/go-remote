FORCE: ;

build: FORCE
	go build -o build/go-remote src/main/main.go

clean:
	rm -rf build
