.PHONEY: test
test: deps lint

	go test -v ./...

.PHONEY: lint
lint: deps

	go vet ./...
	golint -set_exit_status

.PHONEY: deps
deps:

	GO111MODULE=off go get github.com/golang/lint/golint
	go get ./...

.PHONEY: clean
clean:

	go clean
