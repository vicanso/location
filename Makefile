.PHONY: default test test-cover dev

# for dev
dev:
	fresh

# for test
test:
	go test -race -cover ./...

test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out

bench:
	go test -bench=. ./...

build-test:
	go run tool/main.go -max 200
