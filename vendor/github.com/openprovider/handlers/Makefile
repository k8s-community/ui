all: test

check:
	gometalinter ./... --severity=error:vet

test:
	go test -v ./...

prepare_deps:
	go get -u -v -d github.com/takama/router
	go test -i -v ./...

prepare_linters:
	go get -u -v -d github.com/alecthomas/gometalinter
	go install -v github.com/alecthomas/gometalinter
	gometalinter --install --update
