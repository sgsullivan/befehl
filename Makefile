GO_LINKER_SYMBOL := "main.version"

all: test build

clean:
	rm -rf _exe/

build:
	godep go build -x -o _exe/befehl github.com/sgsullivan/befehl/cmd/befehl

test:
	godep go test -v

install: ldflags
	go install -a ${LDFLAGS} ./...

update-deps: godep
	godep save -r ./...

godep:
	go get -u github.com/tools/godep

glv:
	$(eval GO_LINKER_VALUE := $(shell git describe --tags --always))

ldflags: glv
	$(eval LDFLAGS := -ldflags "-X ${GO_LINKER_SYMBOL} ${GO_LINKER_VALUE}")

ver:
	$(eval VERSION := $(shell echo ${GO_LINKER_VALUE} | sed s/^v//))
