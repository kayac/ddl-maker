test: deps lint
	go test -v ./...

lint: deps
	go vet ./...
	go list ./... | xargs -n 1 golint -set_exit_status

deps:
	go get golang.org/x/lint/golint
	go mod download
