all: build

format:
	gofmt -s -w .

clean: clean-test-cache
	rm -rf _exe/
	go clean -x

clean-test-cache:
	go clean -x -testcache

build: format unit-test build-only

build-only:
	go build -x -o _exe/befehl cmd/main/main.go

unit-test: clean-test-cache
	scripts/test/run_unit_tests

integration-test: clean-test-cache
	scripts/test/run_integration_tests

test: unit-test integration-test

update-deps:
	go get -u

release-build:
	scripts/build/release-build

integration-nuke-sshd-hosts:
	integration_tests/nuke_sshd_hosts

integration-start-sshd-hosts:
	integration_tests/start_sshd_hosts 5
