all: test build

clean:
	rm -rf _exe/

build:
	go build -x -o _exe/befehl cmd/main/main.go

test:
	go test -v

update-deps:
	go get -u
