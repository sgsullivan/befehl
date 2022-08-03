all: build

format:
	gofmt -s -w .

clean:
	rm -rf _exe/

build:  format test
	go build -x -o _exe/befehl cmd/main/main.go

test:
	go test -v ./...

update-deps:
	go get -u

release-build:
	scripts/build/release-build
