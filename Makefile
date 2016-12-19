TEST_FILE = $(shell glide novendor)
VET_FILE = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

test: deps lint

	go test -v ${TEST_FILE}

lint: deps

	go tool vet -all -printfuncs=Criticalf,Infof,Warningf,Debugf,Tracef ${VET_FILE}
	glide novendor | xargs -n 1 golint -set_exit_status

deps:

	go get github.com/golang/lint/golint
	glide install

deps-update: clean

	glide update

clean:

	go clean
	glide cache-clear
	rm -rf ./vendor
